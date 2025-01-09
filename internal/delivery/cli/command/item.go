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
		"view":   "item view -i <id> [--done=true] [--all=true] [--tag-names <tagnames>]",
		"delete": "item delete -i <id> [--tags] [--tag-names]",
		"update": "item update -i <id> [-n <name>] [-d <description>] [-t <tag1,tag2>] [--append=true]",
		"done":   "item done -i <id>",
	}

	return &ItemCommand{
		BaseCommand: BaseCommand{
			name:        action,
			description: descriptions[action],
			usage:       usages[action],
			flags:       &Flags{},
		},
		controller: ctrl,
	}
}

func (icmd *ItemCommand) Exec(args []string) error {

	if err := icmd.parseFlags(args); err != nil {
		return fmt.Errorf("can't parse flags : %v", err)
	}
	switch icmd.name {
	case "add":
		return icmd.handleAdd()
	case "view":
		return icmd.handleView()
	case "delete":
		return icmd.handleDelete()
	case "update":
		return icmd.handleUpdate()
	case "done":
		return icmd.handleDone()
	default:
		return flag.ErrHelp
	}
}

func (icmd *ItemCommand) parseFlags(args []string) error {
	fs := flag.NewFlagSet(args[1], flag.ExitOnError)
	resource := args[0]

	fs.StringVar(&icmd.flags.Name, "n", "", fmt.Sprintf("%s name", resource))
	fs.StringVar(&icmd.flags.ID, "i", "", "item id")
	fs.StringVar(&icmd.flags.Description, "d", "", "item description")
	fs.StringVar(&icmd.flags.Tags, "t", "", "item tags")
	fs.StringVar(&icmd.flags.Short, "short", "", "item short view (no tags)")
	fs.StringVar(&icmd.flags.TagNames, "tag-names", "", "item tag names")
	fs.BoolVar(&icmd.flags.All, "all", false, fmt.Sprintf("%s all references (bulk)", resource))
	fs.BoolVar(&icmd.flags.Done, "done", false, "change status of an item to done")

	err := fs.Parse(args[2:])
	return err
}

func (icmd *ItemCommand) handleAdd() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr(icmd.flags.Name, icmd.flags.Description); err != nil {
		return fmt.Errorf("%v", err)
	}
	var tags []string

	if isFlagDefined(icmd.flags.Tags) {
		tags = strings.Split(icmd.flags.Tags, ",")
	}
	err = icmd.controller.AddItem(&model.Item{
		Name:        icmd.flags.Name,
		Description: icmd.flags.Description,
	}, tags...)

	return err

}

// TODO: add single item done view
func (icmd *ItemCommand) handleView() error {
	var err error
	if icmd.flags.All{
		var items []model.Item
		if icmd.flags.Done {
			items, err = icmd.controller.ViewItemsDone()
			if err != nil {
				return err
			}
		} else {
			items, err = icmd.controller.ViewItems()
			if err != nil {
				return err
			}
		}
		response.TabWriter(items)
		return nil
	}

	if err = validation.ValidateFlagsDefinedStr(icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	var item *model.Item
	item, err = icmd.controller.ViewItem(&model.Item{
		ID: icmd.flags.ID,
	})
	if err != nil {
		return err
	}
	response.TabWriter(item)

	return nil

}

func (icmd *ItemCommand) handleUpdate() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr(icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	

	err = icmd.controller.UpdateItem(&model.Item{
		ID: icmd.flags.ID,
	}, map[string]interface{}{
		"name":        icmd.flags.Name,
		"description": icmd.flags.Description,
	})

	return err
}

func (icmd *ItemCommand) handleDone() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr(icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	err = icmd.controller.UpdateItemDone(icmd.flags.ID)

	return err
}

func (icmd *ItemCommand) handleDelete() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr(icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	err = icmd.controller.DeleteItem(&model.Item{
		ID: icmd.flags.ID,
	})

	return err
}
