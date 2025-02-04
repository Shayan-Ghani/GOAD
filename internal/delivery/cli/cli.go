package cli

import (
	"fmt"
	"gocasts/ToDoApp/internal/repository"
	"gocasts/ToDoApp/internal/delivery/cli/command"
	"gocasts/ToDoApp/pkg/validation"
)

type CLI struct {
	repo    repository.Repository
	commands      map[string]command.Command
	validCommands map[string][]string
}

func NewCLI(repo repository.Repository) *CLI {
	return &CLI{
		repo: repo,
		commands:   make(map[string]command.Command),
		validCommands: map[string][]string{
			"item": {"add", "view", "delete", "update", "done", "--help"},
			"tag":  {"view", "delete", "--help"},
		},
	}
}

func Ishelp(arg string) bool {
	if arg == "--help" || arg == "-h" {
		return true
	}
	return false
}

func (c *CLI) registerCommands(resource string, action string) {

	switch resource {
	case "item":
		c.commands[resource] = command.NewItemCommand(c.repo, action)
	case "tag":
		c.commands[resource] = command.NewTagCommand(c.repo, action)
	}

}

func (c *CLI) Exec(args []string) error {
	if err := validation.ValidateFlagCount(len(args), 1); err != nil {
		c.PrintUsage()
		return fmt.Errorf("%w", err)
	}

	if Ishelp(args[0]) {
		c.PrintUsage()
		return nil
	}

	resource := args[0]
	action := args[1]

	validArgcount := 3
	if action == "update" || (action == "add" && resource == "tag") {
		validArgcount = 4
	}

	if err := validation.ValidateResourceAction(c.validCommands, resource, action); err != nil {
		c.PrintUsage()
		return fmt.Errorf("%w", err)
	}

	if !Ishelp(action) {
		if err := validation.ValidateFlagCount(len(args), validArgcount); err != nil {
			c.PrintUsage()
			return fmt.Errorf("%s %s command : %w", resource, action, err)
		}
	}else{
		c.PrintUsage()
		return nil
	}

	c.registerCommands(resource, action)
	return c.commands[resource].Exec(args)
}

func (c *CLI) PrintUsage() {
	fmt.Printf("Available commands:\n")

	itemDesc := map[string]string{
		"--help": "see help for flags and options",
		"add":    "Add a new item",
		"view":   "View an item(s), also use -t <tag-name> instead of -i(single)/--all to see items filtered by that tag name.",
		"delete": "Delete an item or its tags with --del-tags",
		"update": "Update an item",
		"done":   "update item status to done from pending",
	}

	itemUse := map[string]string{
		"--help": "see help for flags and options",
		"add":    "item add -n <name> -d <description> [-t tag1,tag2] [-due-date <date string> (e.g '2025-03-05 15:05:10') ] ",
		"view":   "item view [-i <id>] [--done=true] [--all=true] [-t <items-with-these-tags,tag2>] [--format=json/table]",
		"delete": "item delete -i <id> [-t <tags-to-delete> ] [--del-tags=true]",
		"update": "item update -i <id> [-n <name>] [-d <description>] [-t <tag1,tag2>] [-due-date <date string> (e.g '2025-03-05 15:05:10') ]",
		"done":   "item done -i <id>",
	}

	tagDesc := map[string]string{
		"view":   "View tags",
		"delete": "Delete a tag, all refrences of the tag will be removed from items.",
	}

	tagUse := map[string]string{
		"view":   "tag view --all=true",
		"delete": "tag delete -n <name>",
	}

	fmt.Println("item : ...")
	for act, desc := range itemDesc {
		fmt.Printf("  %-10s\t%s\n", act, desc)
		fmt.Printf("    Usage: %s\n\n", itemUse[act])
	}

	fmt.Printf("\ntag : ...\n")
	for act, desc := range tagDesc {
		fmt.Printf("  %-10s\t%s\n", act, desc)
		fmt.Printf("    Usage: %s\n\n", tagUse[act])
	}
}
