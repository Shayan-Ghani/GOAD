package command

import (
	"flag"
	"fmt"
	"gocasts/ToDoApp/internal/controller"
	"gocasts/ToDoApp/internal/model"
	"gocasts/ToDoApp/internal/pkg/response"
	"gocasts/ToDoApp/internal/pkg/validation"
	"strings"
)

type ItemCommand struct {
	BaseCommand
	controller *controller.SQLTodoController
	action     string
	flags      map[string]interface{}
}

func NewItemCommand(ctrl *controller.SQLTodoController, action string) *ItemCommand {

	descriptions := map[string]string{
		"add":    "Add a new todo item",
		"view":   "View a todo item(s)",
		"delete": "Delete a todo item",
		"update": "Update a todo item",
		"done":   "update item status to done from pending",
	}

	usages := map[string]string{
		"add":    "item add -n <name> -d <description> [-t tag1,tag2]",
		"view":   "item view -i <id> [-done] [--all] [--tag-names <tagnames>]",
		"delete": "item delete -i <id> [--tags] [--tag-names]",
		"update": "item update -i <id> [-n <name>] [-d <description>] [-t/--tags]",
		"done":   "item done <id>",
	}

	return &ItemCommand{
		BaseCommand: BaseCommand{
			name:        action,
			description: descriptions[action],
			usage:       usages[action],
		},
		controller: ctrl,
		action:     action,
	}
}

func (icmd *ItemCommand) Exec(args []string) error {

	if err := icmd.parseFlags(args); err != nil {
		return fmt.Errorf("can't parse flags : %v", err)
	}
	switch icmd.action {
	case "add":
		return icmd.handleAdd()
	case "view":
		return icmd.handleView()
	case "delete":
		return icmd.handleDelete()
	case "update":
		return icmd.handleUpdate()
	default:
		return flag.ErrHelp
	}
}

// TODO: fix all flag type
func (icmd *ItemCommand) parseFlags(args []string) error {
	fs := flag.NewFlagSet(args[1], flag.ExitOnError)
	resource := args[0]
	flags := make(map[string]interface{})
	flags["name"] = *fs.String("n", "", fmt.Sprintf("%s name", resource))
	flags["all"] = *fs.String("all", "", fmt.Sprintf("%s all refrences (bulk)", resource))
	flags["id"] = *fs.String("i", "", "item id")
	flags["description"] = *fs.String("d", "", "item description")
	flags["tags"] = *fs.String("t", "", "item tags")
	flags["short"] = *fs.String("short", "", "item short view (no tags)")
	flags["tagNames"] = *fs.String("tag-names", "", "item tag names")
	flags["done"] = *fs.Bool("done", false, "change status of an item to done")

	err := fs.Parse(args[2:])
	if err != nil {
		icmd.flags = flags
	}
	return err
}

func (icmd *ItemCommand) handleAdd() error {
	if err := validation.ValidateFlagsDefinedStr(icmd.flags, "name", "description"); err != nil {
		return fmt.Errorf("validation failed! %w", err)
	}
	var tags []string

	if icmd.flags["tags"] != "" {
		tags = strings.Split(icmd.flags["tags"].(string), ",")
	}
	err := icmd.controller.AddItem(&model.Item{
		Name:        icmd.flags["name"].(string),
		Description: icmd.flags["description"].(string),
	}, tags...)

	return fmt.Errorf("handleAdd: %s", err)

}

// TODO: add single item done view
func (icmd *ItemCommand) handleView() error {
	var err error
	if icmd.flags["id"] != "" {
		var items []model.Item
		if icmd.flags["done"].(bool) {
			items, err = icmd.controller.ViewItemsDone()
			if err != nil {
				return err
			}
		}else{
			items, err = icmd.controller.ViewItems()
			if err != nil {
				return err
			}
		}
		response.TabWriter(items)
	}

	var item *model.Item
	item, err = icmd.controller.ViewItem(&model.Item{
		ID: icmd.flags["id"].(string),
	})
	if err != nil {
		return err
	}
	response.TabWriter(item)
	
	return nil
	
}


func (icmd *ItemCommand) handleUpdate() error {
	if err := validation.ValidateFlagsDefinedStr(icmd.flags, "id"); err != nil {
		return fmt.Errorf("validation failed! %w", err)
	}

	err := icmd.controller.UpdateItem(&model.Item{
		ID: icmd.flags["id"].(string),
	}, map[string]interface{}{
		"name":        icmd.flags["name"].(string),
		"description": icmd.flags["description"].(string),
	})

	return err
}
func (icmd *ItemCommand) handleDelete() error {
	if err := validation.ValidateFlagsDefinedStr(icmd.flags, "id"); err != nil {
		return fmt.Errorf("validation failed! %w", err)
	}

	err := icmd.controller.DeleteItem(&model.Item{
		ID: icmd.flags["id"].(string),
	})

	return err
}
