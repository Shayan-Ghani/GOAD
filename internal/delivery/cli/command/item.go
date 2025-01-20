package command

import (
	"flag"
	"fmt"
	"gocasts/ToDoApp/internal/controller"
	"gocasts/ToDoApp/internal/pkg/formatter"
	"gocasts/ToDoApp/internal/pkg/response"
	"gocasts/ToDoApp/internal/pkg/validation"
)

type ItemCommand struct {
	BaseCommand
	controller *controller.SQLTodoController
}

func NewItemCommand(ctrl *controller.SQLTodoController, action string) *ItemCommand {
	return &ItemCommand{
		BaseCommand: BaseCommand{
			name:  action,
			flags: &Flags{},
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
		return nil
	}
}

func (icmd *ItemCommand) parseFlags(args []string) error {
	resource := args[0]
	fs := flag.NewFlagSet(resource, flag.ExitOnError)

	fs.StringVar(&icmd.flags.Name, "n", "", fmt.Sprintf("%s name", resource))
	fs.StringVar(&icmd.flags.ID, "i", "", "item id")
	fs.StringVar(&icmd.flags.Description, "d", "", "item description")
	fs.StringVar(&icmd.flags.Tags, "t", "", "item tags to add/delete or filter view by.")
	fs.StringVar(&icmd.flags.Short, "short", "", "item short view (no tags)")
	fs.BoolVar(&icmd.flags.All, "all", false, fmt.Sprintf("when set to true, view all %s references (bulk)", resource))
	fs.BoolVar(&icmd.flags.Done, "done", false, "when set to true, change the status of an item to done")
	fs.BoolVar(&icmd.flags.DelTags, "del-tags", false, "when set to ture, deletes all tags of the item")

	err := fs.Parse(args[2:])
	if args[1] == "--help" {
		fs.PrintDefaults()
	}
	return err
}

func (icmd *ItemCommand) handleAdd() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr([]string{"-n", "-d"}, icmd.flags.Name, icmd.flags.Description); err != nil {
		return fmt.Errorf("%v", err)
	}
	var tags []string

	if isFlagDefined(icmd.flags.Tags) {
		tags = formatter.SplitTags(icmd.flags.Tags)
	}
	err = icmd.controller.AddItem(icmd.flags.Name, icmd.flags.Description, tags...)

	return err

}



func (icmd *ItemCommand) handleView() error {
	handleItems := func(items interface{}, err error) error {
		if err != nil {
			return err
		}
		response.TabWriter(items)
		return nil
	}

	if icmd.flags.Done {
		return handleItems(icmd.controller.ViewItemsDone())
	}
	if icmd.flags.All {
		return handleItems(icmd.controller.ViewItems())
	}
	
	if isFlagDefined(icmd.flags.Tags) {
		tags := formatter.SplitTags(icmd.flags.Tags)
		return handleItems(icmd.controller.ViewItemsByTag(tags))
	}

	if err := validation.ValidateFlagsDefinedStr([]string{"-i"}, icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	return handleItems(icmd.controller.ViewItem(icmd.flags.ID))
}

func (icmd *ItemCommand) handleUpdate() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr([]string{"-i"}, icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	if isFlagDefined(icmd.flags.Tags) {
		icmd.controller.AddItemTag(icmd.flags.ID, formatter.SplitTags(icmd.flags.Tags))
	}

	if !isFlagDefined(icmd.flags.Description, icmd.flags.Name) {
		return fmt.Errorf("description and name must be defined")
	}

	updates := make(map[string]interface{}, 2)

	flagUpdates := map[string]string{
		"name":        icmd.flags.Name,
		"description": icmd.flags.Description,
	}

	for key, value := range flagUpdates {
		if isFlagDefined(value) {
			updates[key] = value
		}
	}

	err = icmd.controller.UpdateItem(icmd.flags.ID, updates)

	return err
}

func (icmd *ItemCommand) handleDone() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr([]string{"-i"}, icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	err = icmd.controller.UpdateItemDone(icmd.flags.ID)

	return err
}

func (icmd *ItemCommand) handleDelete() error {
	var err error

	if err = validation.ValidateFlagsDefinedStr([]string{"-i"}, icmd.flags.ID); err != nil {
		return fmt.Errorf("%w", err)
	}

	if isFlagDefined(icmd.flags.Tags) && icmd.flags.DelTags {
		return validation.New("argument", "Can't use a combinatio of --del-tags and -t")
	}

	if isFlagDefined(icmd.flags.Tags) {
		err = icmd.controller.DeleteItemTags(icmd.flags.ID, formatter.SplitTags(icmd.flags.Tags))
		return err
	}

	if icmd.flags.DelTags {
		err = icmd.controller.DeleteAllItemTags(icmd.flags.ID)
		return err
	}

	err = icmd.controller.DeleteItem(icmd.flags.ID)

	return err
}
