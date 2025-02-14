package sqlrepository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Shayan-Ghani/GOAD/internal/model"
	now "github.com/Shayan-Ghani/GOAD/pkg/time"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepo(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func scanItem(rows *sql.Rows, item *model.Item, row ...*sql.Row) error {

	var isDone []byte
	var createdAt []uint8
	var modifiedAt []uint8
	var dueDate []uint8

	fields := []interface{}{
		&item.ID,
		&item.Name,
		&item.Description,
		&isDone,
		&dueDate,
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
	if err != nil {
		return fmt.Errorf("parsing CreatedAt: %v", err)
	}

	item.DueDate, err = parseDatetime(dueDate)

	return err
}

func (c *ItemRepository) processItemRows(rows *sql.Rows) (itemRows []model.Item, err error) {
	for rows.Next() {
		var item model.Item
		if err = scanItem(rows, &item); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
	
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

func (c *ItemRepository) AddItem(name string, description string, dueDate time.Time, tags ...string) (insID int64 , err error) {
	q := "INSERT INTO items (name, description, created_at) VALUES (?, ?, ?)"
	args := []interface{}{name, description, now.Now()}

	if !dueDate.IsZero() {
		q = "INSERT INTO items (name, description, created_at, due_date) VALUES (?, ?, ?, ?)"
		args = append(args, dueDate)
	}

	stmt, err := c.db.Prepare(q)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	insert, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert item: %w", err)
	}

	insID, err = insert.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't get the last insert ID: %v", err)
	}
	return insID, nil
}

func (c *ItemRepository) GetItem(id string) (*model.Item, error) {
	q := "SELECT * FROM items WHERE id = ?"
	stmt, err := c.db.Prepare(q)
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

	return item, err
}

type QueryTemplate struct {
	Condition string
	Args      []interface{}
}

func (c *ItemRepository) GetItems(templates ...QueryTemplate) ([]model.Item, error) {
	q := "SELECT * FROM items"
	var args []interface{}

	if len(templates) > 0 {
		template := templates[0]
		if template.Condition != "" {
			q += fmt.Sprintf(" WHERE %s", template.Condition)
			args = template.Args
		}
	}

	var rows *sql.Rows
	var err error
	if len(args) > 0 {
		rows, err = c.db.Query(q, args...)
	} else {
		rows, err = c.db.Query(q)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	return c.processItemRows(rows)
}

func (c *ItemRepository) GetItemsDone() ([]model.Item, error) {
	return c.GetItems(QueryTemplate{
		Condition: "is_done = ?",
		Args:      []interface{}{1},
	})
}

func (c *ItemRepository) GetItemByTag(tags []string) ([]model.Item, error) {
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
		return nil, fmt.Errorf("failed to execute query Get item tags: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()
	return c.processItemRows(rows)
}

func (c *ItemRepository) DeleteItem(id string) error {
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

func (c *ItemRepository) UpdateItem(id string, updates map[string]interface{}) error {

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

func (c *ItemRepository) UpdateItemStatus(id string) error {
	q := "UPDATE items SET is_done = 1 WHERE id = ?"
	stmt, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("couldn't prepare update item is_done : %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
