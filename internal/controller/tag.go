package controller

import (
	"fmt"
	"gocasts/ToDoApp/internal/model"
	now "gocasts/ToDoApp/internal/pkg/time"
	"strings"
)

func (c *SQLTodoController) AddTag(tags []string) error {

	params := make([]interface{}, 0, len(tags)*2)
	placeHolders := make([]string, len(tags))

	for i, tag := range tags {
		placeHolders[i] = "(?,?)"
		params = append(params, tag, now.Now())
	}

	q := fmt.Sprintf("INSERT IGNORE INTO tags (name, created_at) VALUES %s", strings.Join(placeHolders, ","))

	stmtIns, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to prepare add tag statement: %v ", err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(params...)
	return err
}

func (c *SQLTodoController) AddTagInto(tag *model.Tag) error {
	q := "INSERT IGNORE INTO tags (name,created_at) VALUES (?, ?)"
	// a := `
	// INSERT IGNORE INTO tags (name)
	// VALUES
	// ('Electronics'),
	// ('Gadgets'),
	// ('Smartphones');
	// `
	stmtIn, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmtIn.Close()

	_, err = stmtIn.Exec(tag.Name)

	return err
}

func (c *SQLTodoController) DeleteTag(tag *model.Tag) error {
	stmtDel, err := c.db.Prepare("DELETE from tags where name = ?")
	if err == nil {
		return fmt.Errorf("could not prepare delete tag statement: %v", err)
	}
	defer stmtDel.Close()

	if _, err := stmtDel.Exec(tag.Name); err != nil {
		return fmt.Errorf("could not query delete tag statement: %v", err)
	}
	return nil
}