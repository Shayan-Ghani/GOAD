package sqlrepository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Shayan-Ghani/GOAD/internal/model"
	now "github.com/Shayan-Ghani/GOAD/pkg/time"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

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

func (c *TagRepository) AddTag(tags []string) error {
	if err := validation.ValidateTagNames(tags); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	params, placeHolders, _ := packTagParamsAndPlacholders(tags, false, len(tags))

	q := fmt.Sprintf("INSERT IGNORE INTO tags (name, created_at) VALUES %s", placeHolders)

	stmtIns, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to prepare add tag statement: %v ", err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(params...)
	return err
}

func (c *TagRepository) AddTagInto(name string) error {
	q := "INSERT IGNORE INTO tags (name,created_at) VALUES (?, ?)"
	stmtIn, err := c.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmtIn.Close()

	_, err = stmtIn.Exec(name)

	return err
}

func (c *TagRepository) GetTags() ([]model.Tag, error) {
	return nil, nil
}

func (c *TagRepository) getTagID(name string) (string, error) {
	var id string

	stmt, err := c.db.Prepare("SELECT id FROM tags WHERE name = ?")
	if err != nil {
		return id, fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(name).Scan(&id)
	return id, err
}

func (c *TagRepository) DeleteTag(name string) error {
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

func (c *TagRepository) AddItemTag(id string, tags []string) error {
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

func (c *TagRepository) GetItemTagsName(id string) ([]string, error) {
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

func (c *TagRepository) DeleteItemTags(id string, tags []string) error {
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

func (c *TagRepository) DeleteAllItemTags(id string) error {
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

