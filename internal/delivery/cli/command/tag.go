package command

import (
	"fmt"
	"gocasts/ToDoApp/internal/controller"
)

type TagCommand struct {
	BaseCommand
	controller *controller.SQLTodoController
	action     string
}

func NewTagCommand(ctrl *controller.SQLTodoController, action string) *TagCommand {

	descriptions := map[string]string{
		"view":   "View tags",
		"delete": "Delete a tag, all refrences of the tag will be removed from items.",
	}

	usages := map[string]string{
		"view":   "tag view -i <id> [-done] [--all] [-tag <tagname>]",
		"delete": "tag delete -i <id>",
	}

	return &TagCommand{
		BaseCommand: BaseCommand{
			name:        action,
			description: descriptions[action],
			usage:       usages[action],
		},
		controller: ctrl,
		action:     action,
	}
}

func (tcmd *TagCommand) Exec(flags []string) error {

	switch tcmd.action {
	case "add":
		return tcmd.handleAdd(flags)
	case "view":
		return tcmd.handleView(flags)
	case "delete":
		return tcmd.handleDelete(flags)
	default:
		return fmt.Errorf("unknown action (%s) for Item", tcmd.action)
	}
}

func (tcmd *TagCommand) handleAdd(flags []string) error {
	return nil
}
func (tcmd *TagCommand) handleView(flags []string) error {
	return nil
}
func (tcmd *TagCommand) handleDelete(flags []string) error {
	return nil
}
