package controller

import (
	"database/sql"
	"fmt"
	"gocasts/ToDoApp/internal/model"
	now "gocasts/ToDoApp/internal/pkg/time"
	"strconv"
	"strings"
	"time"
)

type SQLTodoController struct {
	db *sql.DB
}

func NewSQLTodoController(db *sql.DB) *SQLTodoController {
	return &SQLTodoController{db: db}
}

func scanItem(rows *sql.Rows, item *model.Item, row ...*sql.Row) error {
	var createdAt []uint8
	var modifiedAt []uint8
	var isDone []byte

	fields := []interface{}{
		&item.ID,
		&item.Name,
		&item.Description,
		&isDone,
		&createdAt,
		&modifiedAt,
	}

	var err error
	if row != nil {
		err = row[0].Scan(fields...)
	} else {
		err = rows.Scan(fields...)
	}

	if err != nil {
		return err
	}

	item.IsDone = len(isDone) > 0 && isDone[0] == 1

	if len(createdAt) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err == nil {
			item.CreatedAt = parsedTime
		}
	}

	if len(modifiedAt) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(modifiedAt))
		if err == nil {
			item.ModifiedAt = parsedTime
		}
	}
	return nil
}

func (c *SQLTodoController) AddItem(item *model.Item, tags ...string) error {
	stmtIn, err := c.db.Prepare("INSERT INTO items (name, description, created_at ) VALUES (?,?,?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmtIn.Close()

	test, err := stmtIn.Exec(item.Name, item.Description, now.Now()); 
	
	if err != nil {
		return fmt.Errorf("failed to insert item: %v", err)
	}
	
	if tags != nil {	
		id , err := test.LastInsertId()
		if err != nil {
			return fmt.Errorf("eRRRRRRRRRR: %v", err)
		}	
		if err = c.AddItemTag(strconv.Itoa(int(id)), tags); err != nil {
			return fmt.Errorf(fmt.Sprintf("fail to add tags %s to item: ", tags), err)
		}

		if err = c.AddTag(tags); err != nil {
			return fmt.Errorf(fmt.Sprintf("fail to add tags %s to tags table: ", tags), err)
		}
	}

	return nil
}

func (c *SQLTodoController) AddItemTag(id string, tags []string) error {
	params, placeHolders, _ := packTagParamsAndPlacholders(tags, true ,len(tags))

	q := fmt.Sprintf(`INSERT INTO item_tags (item_id, tag_id)
SELECT i.id, t.id
FROM items i
JOIN tags t ON t.name IN (%s)
WHERE i.id = ?`, placeHolders )

	stmtIns , err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to Prepare addItemTag statement: %v", err)
	}
	defer stmtIns.Close()
	params = append(params, id)
	_, err = stmtIns.Exec(params...)
	return err
}





func (c *SQLTodoController) ViewItem(item *model.Item) (*model.Item, error) {
	stmt, err := c.db.Prepare("SELECT * FROM items WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	if err := scanItem(nil, item, stmt.QueryRow(item.ID)); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	return item, err
}

// TODO: chekOut for Security of variadic query parameter.
func (c *SQLTodoController) ViewItems(item *model.Item, query ...string) ([]model.Item, error) {
	q := "SELECT * FROM items"
	if query != nil {
		q = query[0]
	}
	stmt, err := c.db.Prepare(q)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var ItemRows []model.Item

	if rows.Next() {
		if err = scanItem(rows, item); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		ItemRows = append(ItemRows, *item)
	}

	return ItemRows, nil
}

func (c *SQLTodoController) ViewItemsDone(item *model.Item) ([]model.Item, error) {
	return c.ViewItems(item, "SELECT * FROM items_done")
}

func (c *SQLTodoController) ViewItemTagsName(item *model.Item) ([]model.Tag, error) {
	q := `SELECT name 
	FROM tags where id in (
	SELECT tag_id
	FROM item_tags WHERE item_id = ?
	)`
	fmt.Println(q, item.ID)
	return nil, nil
}

func (c *SQLTodoController) ViewItemsByTag(tag *model.Tag) ([]model.Item, error) {
	q := `SELECT id, name 
FROM items  
WHERE id in (
	SELECT item_id
	FROM item_tags WHERE tag_id = (
	SELECT id
	FROM tags
	WHERE name = '?'
	)
)`
	fmt.Println(q, tag.Name)
	return nil, nil
}

func (c *SQLTodoController) DeleteItem(item *model.Item) error {
	stmtDel, err := c.db.Prepare("DELETE FROM items WHERE id = ?")
	if err != nil {
		return fmt.Errorf("could not prepare delete item statement: %v", err)
	}

	defer stmtDel.Close()

	if _, err := stmtDel.Exec(item.ID); err != nil {
		return fmt.Errorf("could not query delete item statement: %v", err)
	}
	return nil
}

func (c *SQLTodoController) DeleteItemTag(TagName string) error {
	stmtDel, err := c.db.Prepare("DELETE from item_tags where tag_id = (SELECT id FROM tags where name = ?)")
	if err == nil {
		return fmt.Errorf("could not prepare delete tag statement: %v", err)
	}
	defer stmtDel.Close()

	if _, err := stmtDel.Exec(TagName); err != nil {
		return fmt.Errorf("could not query delete tag statement: %v", err)
	}
	return nil

}



func (c *SQLTodoController) UpdateItem(item *model.Item, updates map[string]interface{}) error {
	// Build query dynamically
	var setFields []string
	var args []interface{}

	for field, value := range updates {
		setFields = append(setFields, field+" = ?")
		args = append(args, value)
	}

	// Add ID to args
	args = append(args, item.ID)

	query := fmt.Sprintf("UPDATE items SET %s WHERE id = ?", strings.Join(setFields, ", "))

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("could not prepare update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	return err
}

func (c *SQLTodoController) UpdateItemDone(id string) error {
	q := "UPDATE items SET is_done = 1 WHERE id = ?"
	stmt,err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("couldn't prepare update item is_done : %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// func (c *SQLTodoController) UpdateItemName(id string, name string) error {

// 	return nil
// }

/// Getter and setters

// func (c *SQLTodoController) GetId() string {
// 	return c.Id
// }

// //	func (c *SQLTodoController) GetName() string {
// //		return c.name
// //	}
// func (c *SQLTodoController) GetTag(tag string) string {
// 	return fmt.Sprintf(tag)
// }
// func (c *SQLTodoController) GetTags() []string {
// 	return c.tag
// }

// // func (c *SQLTodoController) GetDesc() string {
// // 	return c.description
// // }

// func (c *SQLTodoController) SetName(name string) {
// 	c.name = name
// }
// func (c *SQLTodoController) SetTag(tag ...string) {
// 	c.tag = tag
// }

// func (c *SQLTodoController) SetDesc(desc string) {
// 	c.description = desc
// }
