package cli

import (
	"fmt"
	"gocasts/ToDoApp/internal/controller"
	"gocasts/ToDoApp/internal/delivery/cli/command"
	"gocasts/ToDoApp/internal/pkg/validation"
)

type CLI struct {
	controller    *controller.SQLTodoController
	commands      map[string]command.Command
	validCommands map[string][]string
}

func NewCLI(controller *controller.SQLTodoController) *CLI {
	return &CLI{
		controller: controller,
		commands:   make(map[string]command.Command),
		validCommands: map[string][]string{
			"item": {"add", "view", "delete", "update", "done"},
			"tag":  {"add", "view", "delete"},
		},
	}
}

func (c *CLI) registerCommands(resource string, action string) {

	switch resource {
	case "item":
		c.commands[resource] = command.NewItemCommand(c.controller, action)
	case "tag":
		c.commands[resource] = command.NewTagCommand(c.controller, action)
	}
	
}

func (c *CLI) Exec(args []string) error {
	
	if err := validation.ValidateFlagCount(len(args), 3); err != nil {
		c.PrintUsage()
		return fmt.Errorf("validation failed! %w", err)
	} 
	resource := args[0]
	action := args[1]

	validArgcount := 3
	if action == "update" || (action == "add" && resource == "tag") {
		validArgcount = 4
	}
	if err := validation.ValidateResourceAction(c.validCommands, resource, action); err != nil {
		return fmt.Errorf("validation failed! %w", err)
	}
	if err := validation.ValidateFlagCount(len(args), validArgcount); err != nil {
		return fmt.Errorf("validation failed! %s %s command : %w", resource, action, err)
	}

	c.registerCommands(resource, action)

	return c.commands[resource].Exec(args[1:])
}

func (c *CLI) PrintUsage() {
	fmt.Println("Available commands:")
	for rsc, cmd := range c.commands {
		fmt.Printf("\n%s commands:\n", rsc)
		fmt.Printf("  %-10s\t%s\n", cmd.Name(), cmd.Description())
		fmt.Printf("    Usage: %s\n", cmd.Usage())
	}
}
