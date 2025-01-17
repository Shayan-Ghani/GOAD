package command

import (
	"flag"
	"fmt"
	"gocasts/ToDoApp/internal/controller"
	"gocasts/ToDoApp/internal/model"
	"gocasts/ToDoApp/internal/pkg/response"
	"gocasts/ToDoApp/internal/pkg/validation"
)

type TagCommand struct {
	BaseCommand
	controller *controller.SQLTodoController
}

func NewTagCommand(ctrl *controller.SQLTodoController, action string) *TagCommand {
	return &TagCommand{
		BaseCommand: BaseCommand{
			name:        action,
			flags:       &Flags{},
		},
		controller: ctrl,
	}
}

func (tcmd *TagCommand) Exec(args []string) error {
	if err := tcmd.parseFlags(args); err != nil {
		return fmt.Errorf("can't parse flags : %v", err)
	}

	switch tcmd.name {
	case "view":
		return tcmd.handleView()
	case "delete":
		return tcmd.handleDelete()
	default:
		return fmt.Errorf("unknown action (%s) for Item", tcmd.name)
	}
}

func (tcmd *TagCommand) parseFlags(args []string) error {
	resource := args[0]
	fs := flag.NewFlagSet(resource, flag.ExitOnError)

	fs.StringVar(&tcmd.flags.Name, "n", "", fmt.Sprintf("%s name", resource))

	err := fs.Parse(args[2:])
	if args[1] == "--help"{
		fs.PrintDefaults()
	}
	return err
}

func (tcmd *TagCommand) handleView() error {
	tags, err := tcmd.controller.ViewTags()
	if err != nil {
		return err
	}
	response.TabWriter(tags)
	return err
}
func (tcmd *TagCommand) handleDelete() error {
	var err error
	if err = validation.ValidateFlagsDefinedStr([]string{"-n"}, tcmd.flags.Name); err != nil {
		return fmt.Errorf("%w", err)
	}
	err = tcmd.controller.DeleteTag(&model.Tag{
		Name: tcmd.flags.Name,
	})
	return err
}
