package command

import (
	"fmt"
	"gocasts/ToDoApp/internal/controller"
)

type ItemCommand struct {
	BaseCommand
	controller *controller.SQLTodoController
	action     string
}

func NewItemCommand(ctrl *controller.SQLTodoController, action string) *ItemCommand {

	descriptions := map[string]string{
		"add":    "Add a new todo item",
		"view":   "View a todo item(s)",
		"delete": "Delete a todo item",
		"update": "Update a todo item",
		"done": "update item status to done from pending",
	}

	usages := map[string]string{
		"add":    "item add -n <name> -d <description> [-t tag1,tag2]",
		"view":   "item view <id> [-done] [all] [-tag <tagname>]",
		"delete": "item delete <id> [--tags]",
		"update": "item update <id> [-n <name>] [-d <description>] [-t/--tags]",
		"done": "item done <id>",
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

func (icmd *ItemCommand) Exec(flags []string) error {

	switch icmd.action {
	case "add":
		return icmd.handleAdd(flags)
	case "view":
		return icmd.handleView(flags)
	case "delete":
		return icmd.handleDelete(flags)
	case "update":
		return icmd.handleUpdate(flags)
	default:
		return fmt.Errorf("unknown action (%s) for Item", icmd.action)
	}
}

func (icmd *ItemCommand) handleAdd(flags []string) error {
	return nil
}
func (icmd *ItemCommand) handleView(flags []string) error {
	return nil
}
func (icmd *ItemCommand) handleUpdate(flags []string) error {
	return nil
}
func (icmd *ItemCommand) handleDelete(flags []string) error {
	return nil
}
