package controller

import (
	"database/sql"
	"fmt"
	"gocasts/ToDoApp/internal/model"
	now "gocasts/ToDoApp/internal/pkg/time"
	"gocasts/ToDoApp/internal/pkg/validation"
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
	var isDone []byte
	var createdAt []uint8
	var modifiedAt []uint8

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

	item.CreatedAt, err = parseDatetime(createdAt)

	if err != nil {
		return fmt.Errorf("parsing CreatedAt: %v", err)
	}
	item.ModifiedAt, err = parseDatetime(modifiedAt)
	return err
}

func (c *SQLTodoController) processItemRows(rows *sql.Rows) (itemRows []model.Item, err error) {
	for rows.Next() {
		var item model.Item
		if err = scanItem(rows, &item); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		item.TagsNames, err = c.ViewItemTagsName(item.ID)
		if err != nil {
			return nil, fmt.Errorf("getting item tags: %v", err)
		}
		itemRows = append(itemRows, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	return itemRows, err
}

func parseDatetime(datetime []uint8) (time.Time, error) {

	if len(datetime) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(datetime))
		if err != nil {
			return time.Time{}, fmt.Errorf("can't parse datetime (%s): %v", datetime, err)
		}
		return parsedTime, nil
	}
	return time.Time{}, nil
}

func (c *SQLTodoController) AddItem(name string, description string, tags ...string) error {
	stmtIn, err := c.db.Prepare("INSERT INTO items (name, description, created_at ) VALUES (?,?,?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmtIn.Close()

	insert, err := stmtIn.Exec(name, description, now.Now())

	if err != nil {
		return fmt.Errorf("failed to insert item: %v", err)
	}

	if tags != nil {
		id, err := insert.LastInsertId()
		if err != nil {
			return fmt.Errorf("couldn't get the last insert ID: %v", err)
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
	if err := validation.ValidateTagNames(tags); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	var tagsToAdd []string
	for _, tag := range tags {
		if _, err := c.getTagID(tag); err != nil {
			if err == sql.ErrNoRows {
				tagsToAdd = append(tagsToAdd, tag)
			} else {
				return err
			}
		}
	}
	if tagsToAdd != nil {
		if err := c.AddTag(tagsToAdd); err != nil {
			return err
		}
	}

	params, placeHolders, _ := packTagParamsAndPlacholders(tags, true, len(tags))

	q := fmt.Sprintf(`INSERT IGNORE INTO item_tags (item_id, tag_id)
SELECT i.id, t.id
FROM items i
JOIN tags t ON t.name IN (%s)
WHERE i.id = ?`, placeHolders)

	stmtIns, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to Prepare addItemTag statement: %v", err)
	}
	defer stmtIns.Close()
	params = append(params, id)
	_, err = stmtIns.Exec(params...)
	return err
}

func (c *SQLTodoController) ViewItem(id string) (*model.Item, error) {
	stmt, err := c.db.Prepare("SELECT * FROM items WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	item := &model.Item{}
	if err := scanItem(nil, item, stmt.QueryRow(id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	item.TagsNames, err = c.ViewItemTagsName(id)
	if err != nil {
		return nil, fmt.Errorf("getting item tags: %v", err)
	}

	return item, err
}

// TODO: chekOut for Security of variadic query parameter.
func (c *SQLTodoController) ViewItems(query ...string) ([]model.Item, error) {
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

	rws, err := c.processItemRows(rows)
	return rws, err
}

func (c *SQLTodoController) ViewItemsDone() ([]model.Item, error) {
	return c.ViewItems("SELECT * FROM items WHERE is_done = 1")
}

func (c *SQLTodoController) ViewItemTagsName(id string) ([]string, error) {
	q := `SELECT name 
	FROM tags where id in (
	SELECT tag_id
	FROM item_tags WHERE item_id = ?
	)`

	stmt, err := c.db.Prepare(q)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare item tag name statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query item tag name: %v", err)
	}

	defer rows.Close()

	var itemTags []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("can't read item tag name row: %v", err)
		}
		itemTags = append(itemTags, name)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return itemTags, nil
}

func (c *SQLTodoController) ViewItemsByTag(tags []string) ([]model.Item, error) {
	params, placeHolders, _ := packTagParamsAndPlacholders(tags, true, len(tags))
	q := fmt.Sprintf(`SELECT *
FROM items  
WHERE id IN (
	SELECT item_id
	FROM item_tags WHERE tag_id IN (
	SELECT id
	FROM tags
	WHERE name IN (%s)
	)
)`, placeHolders)
	stmt, err := c.db.Prepare(q)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query view item tags: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	rws, err := c.processItemRows(rows)
	return rws, err
}

func (c *SQLTodoController) DeleteItem(id string) error {
	stmtDel, err := c.db.Prepare("DELETE FROM items WHERE id = ?")
	if err != nil {
		return fmt.Errorf("could not prepare delete item statement: %v", err)
	}

	defer stmtDel.Close()

	if _, err := stmtDel.Exec(id); err != nil {
		return fmt.Errorf("could not query delete item statement: %v", err)
	}
	return nil
}

func (c *SQLTodoController) DeleteItemTags(id string, tags []string) error {
	args := make([]interface{}, 0)
	args = append(args, id)

	params, placeHolders, _ := packTagParamsAndPlacholders(tags, true, len(tags))
	q := fmt.Sprintf("DELETE from item_tags where item_id = ? and tag_id IN (SELECT id FROM tags where name IN (%s))", placeHolders)
	stmtDel, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to Prepare DeleteItemTags statement: %v", err)
	}

	defer stmtDel.Close()

	args = append(args, params...)
	if _, err := stmtDel.Exec(args...); err != nil {
		return fmt.Errorf("could not query delete tag statement: %v", err)
	}
	return nil
}

func (c *SQLTodoController) DeleteAllItemTags(id string) error {
	q := "DELETE from item_tags where item_id = ?"
	stmtDel, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to Prepare DeleteAllItemTags statement: %v", err)
	}

	defer stmtDel.Close()
	if _, err := stmtDel.Exec(id); err != nil {
		return fmt.Errorf("could not query all delete tag statement: %v", err)
	}
	return nil
}

func (c *SQLTodoController) UpdateItem(id string, updates map[string]interface{}) error {

	var setFields []string
	var args []interface{}

	for field, value := range updates {
		setFields = append(setFields, field+" = ?")
		args = append(args, value)
	}

	args = append(args, id)

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
	stmt, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("couldn't prepare update item is_done : %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
