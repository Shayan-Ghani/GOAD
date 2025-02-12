package cli

import (
	"fmt"

	"github.com/Shayan-Ghani/GOAD/internal/model"
	"github.com/Shayan-Ghani/GOAD/pkg/response"
)

type Printer struct{}

func NewPrinter() Printer {
	return Printer{}
}

func (p Printer) PrintResponse(format string, items []model.Item) {
	response.Respond(format, items)
}

func (p Printer) PrintUsage(s ...string) {

	if s[0] != "" {
		fmt.Println(s[0])
		return
	}

	fmt.Println("Available commands:")

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
