package controller

import (
	"fmt"
	"gocasts/ToDoApp/internal/model"
	now "gocasts/ToDoApp/internal/pkg/time"
	"gocasts/ToDoApp/internal/pkg/validation"
	"strings"
)

func packTagParamsAndPlacholders(tags []string, itemTag bool, tagCount int) ([]interface{}, string, []string) {
	params := make([]interface{}, 0, tagCount*2)
	placeHolders := make([]string, tagCount)

	for i, tag := range tags {
		if !itemTag {
			placeHolders[i] = "(?,?)"
			params = append(params, tag, now.Now())
		} else {
			placeHolders[i] = "?"
			params = append(params, tag)
		}
	}

	return params, strings.Join(placeHolders, ","), placeHolders
}

func (c *SQLTodoController) AddTag(tags []string) error {
	if err := validation.ValidateTagNames(tags); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	params, placeHolders, _ := packTagParamsAndPlacholders(tags, false, len(tags))

	q := fmt.Sprintf("INSERT INTO tags (name, created_at) VALUES %s", placeHolders)

	stmtIns, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to prepare add tag statement: %v ", err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(params...)
	return err
}

func (c *SQLTodoController) AddTagInto(name string) error {
	q := "INSERT IGNORE INTO tags (name,created_at) VALUES (?, ?)"
	stmtIn, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmtIn.Close()

	_, err = stmtIn.Exec(name)

	return err
}

func (c *SQLTodoController) ViewTags() ([]model.Tag, error) {
	return nil, nil
}

func (c *SQLTodoController) getTagID(name string) (string, error) {
	var id string

	stmt, err := c.db.Prepare("SELECT id FROM tags WHERE name = ?")
	if err != nil {
		return id, fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(name).Scan(&id)
	return id, err
}

func (c *SQLTodoController) DeleteTag(name string) error {
	stmtDel, err := c.db.Prepare("DELETE from tags where name = ?")
	if err != nil {
		return fmt.Errorf("could not prepare delete tag statement: %v", err)
	}
	defer stmtDel.Close()

	if _, err := stmtDel.Exec(name); err != nil {
		return fmt.Errorf("could not query delete tag statement: %v", err)
	}
	return err
}
