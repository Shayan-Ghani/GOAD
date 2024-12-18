package controller

import (
	"database/sql"
	"fmt"
	"gocasts/ToDoApp/internal/model"
	now "gocasts/ToDoApp/internal/pkg/time"
	"strings"
)

type SQLTodoController struct {
	db *sql.DB
}

func NewSQLTodoController(db *sql.DB) *SQLTodoController {
	return &SQLTodoController{db: db}
}

func scanItem(rows *sql.Rows, item *model.Item, row ...*sql.Row) error {
	fields := []interface{}{
		&item.ID,
		&item.Name,
		&item.Description,
		&item.CreatedAt,
		&item.ModifiedAt,
		&item.DeletedAt,
		&item.Tags,
	}

	if row != nil {
		return row[0].Scan(fields...)
	}
	return rows.Scan(fields...)
}

func (c *SQLTodoController) AddItem(item *model.Item) error {
	stmtIn, err := c.db.Prepare("INSERT INTO items (name, description, created_at ) VALUES (?,?,?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmtIn.Close()

	if _, err = stmtIn.Exec(item.Name, item.Description, now.Now()); err != nil {
		return fmt.Errorf("failed to insert item: %v", err)
	}

	return nil
}

func (c *SQLTodoController) ViewItem(item *model.Item) (*model.Item, error) {
	stmt, err := c.db.Prepare("SELECT * FROM items WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	if err := scanItem(nil,item, stmt.QueryRow(item.ID)); err != nil {
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

// func (c *SQLTodoController) ViewItemsByTag(tag *model.Tag) ([]model.Item, error) {
// 	return nil, nil
// }

// func (c *SQLTodoController) GetItemByTag(tag string) (error , string) {

// 	// if !TagExists(tag){
// 	// 	return fmt.Errorf("tag doesn't exist"), ""
// 	// }
// 	/// join to get the tag name and item
// 	// err, tag_id = sql.Conn("select * from tags where name == ?")
// 	// err, item = sql.Conn("Select * from items where tag_id == ?")
// }

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

func (c *SQLTodoController) DeleteTag(tag *model.Tag) error {
	stmtDel, err := c.db.Prepare("DELETE from tags where name = ?")
	if err == nil {
		return fmt.Errorf("could not prepare delete tag statement: %v", err)
	}
	defer stmtDel.Close()

	if _, err := stmtDel.Exec(tag.ID); err != nil {
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

// func (c *SQLTodoController) UpdateItemName(id uint, name string) error {

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
