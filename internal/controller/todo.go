package controller

import (
	"database/sql"
	"fmt"
	now "gocasts/ToDoApp/internal/pkg/time"
	"strings"
)

type SQLTodoController struct {
	db *sql.DB
}

func NewSQLTodoController(db *sql.DB) *SQLTodoController {
	return &SQLTodoController{db: db}
}

func (c *SQLTodoController) AddItem(name string, desc string, tags ...string) error {
	// var itemTags []model.Tag
	// for _, tagName := range tags{
	// 	model.Tag.Name=
	// }

	stmtIn, err := c.db.Prepare("INSERT INTO items (name, description, created_at ) VALUES (?,?,?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmtIn.Close()

	if _, err = stmtIn.Exec(name, desc, now.Now()); err != nil {
		return fmt.Errorf("failed to insert item: %v", err)
	}

	return nil
}

func (c *SQLTodoController) ViewItem(id uint) (string, error) {
	stmt , err := c.db.Prepare("SELECT * FROM items WHERE id = ?")
	if err != nil {
		return "", fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

    var item string
	
	if err := stmt.QueryRow(id).Scan(&item); err != nil {
		if err == sql.ErrNoRows {
            return "", fmt.Errorf("item not found")
        }
		return "", fmt.Errorf("failed to scan row: %v", err)
	}
	
	return item , nil
}

func (c *SQLTodoController) ViewItems(id uint) (string, error) {
    stmt, err := c.db.Prepare("SELECT * FROM items")
    if err != nil {
        return "", fmt.Errorf("failed to prepare statement: %v", err)
    }
    defer stmt.Close() 

    rows, err := stmt.Query() 
    if err != nil {
        return "", fmt.Errorf("failed to execute query: %v", err)
    }
    defer rows.Close()

    var item string
    if rows.Next() { 
        if err = rows.Scan(&item); err != nil {
            return "", fmt.Errorf("failed to scan row: %v", err)
        }
    }

    return item, nil
}

func (c *SQLTodoController) ViewItemsDone(id string) (string, error) {
	stmt, err := c.db.Prepare("SELECT * FROM items_done")
    if err != nil {
        return "", fmt.Errorf("failed to prepare statement: %v", err)
    }
    defer stmt.Close() 

    rows, err := stmt.Query() 
    if err != nil {
        return "", fmt.Errorf("failed to execute query: %v", err)
    }
    defer rows.Close()

    var item string
    if rows.Next() { 
        if err = rows.Scan(&item); err != nil {
            return "", fmt.Errorf("failed to scan row: %v", err)
        }
    }

    return item, nil
}

// func (c *SQLTodoController) ViewItemByTag(tag string) string {
// 	return c.GetItemByTag(tag)
// }

// func (c *SQLTodoController) GetItemByTag(tag string) (error , string) {

// 	// if !TagExists(tag){
// 	// 	return fmt.Errorf("tag doesn't exist"), ""
// 	// }
// 	/// join to get the tag name and item
// 	// err, tag_id = sql.Conn("select * from tags where name == ?")
// 	// err, item = sql.Conn("Select * from items where tag_id == ?")
// }

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

func (c *SQLTodoController) DeleteTag(tag string) error {
	stmtDel, err := c.db.Prepare("DELETE from tags where name = ?")
	if err == nil {
		return fmt.Errorf("could not prepare delete tag statement: %v", err)
	}
	defer stmtDel.Close()

	if _ , err := stmtDel.Exec(tag); err != nil {
		return fmt.Errorf("could not query delete tag statement: %v", err)
	}
	return nil
}

func (c *SQLTodoController) UpdateItem(id uint, updates map[string]interface{}) error {
	// Build query dynamically
	var setFields []string
	var args []interface{}
	
	for field, value := range updates {
		setFields = append(setFields, field + " = ?")
		args = append(args, value)
	}
	
	// Add ID to args
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
