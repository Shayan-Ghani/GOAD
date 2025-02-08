package cmdflag

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

func IsFlagDefined(flags ...string) bool {
	for _, f := range flags {
		if f == "" {
			return false
		}
	}
	return true
}

type Flags struct {
	resource    string
	command     string
	Name        string
	ID          string
	Description string
	Tags        string
	Short       string
	Format      string
	DueDate     string
	All         bool
	Done        bool
	DelTags     bool
}

func New(resource string, command string) *Flags {
	return &Flags{
		resource: resource,
		command: command,
	}
}

func (f *Flags) Parse(args []string) error {
	helpStr := fmt.Sprintf("Available Flags for %s %s: \n\n", f.resource, f.command)
	buff := bytes.NewBufferString(helpStr)
	
	fs := flag.NewFlagSet(f.resource, flag.ExitOnError)
	fs.SetOutput(buff)
	
	fs.StringVar(&f.Name, "n", "", fmt.Sprintf("%s name", f.resource))
	fs.StringVar(&f.ID, "i", "", "item id")
	fs.StringVar(&f.Description, "d", "", "item description")
	fs.StringVar(&f.Tags, "t", "", "item tags to add/delete or filter view by.")
	fs.StringVar(&f.Short, "short", "", "item short view (no tags)")
	fs.StringVar(&f.Format, "format", "table", "output format")
	fs.StringVar(&f.DueDate, "due-date", "", "item due date in 'year-month-day hour:minute:second' format.")
	fs.BoolVar(&f.All, "all", false, fmt.Sprintf("when set to true, view all %s references (bulk)", f.resource))
	fs.BoolVar(&f.Done, "done", false, "when set to true, change the status of an item to done")
	fs.BoolVar(&f.DelTags, "del-tags", false, "when set to ture, deletes all tags of the item")



	if args[1] == "--help" {
		fs.PrintDefaults()
		return validation.Help{
			Message : buff.String(),
		}
	}

	if err := fs.Parse(args[2:]); err != nil{
		return validation.Help{
			Message : err.Error(),
		}
	}
	
	
	return nil

}


func (f *Flags) CheckCommandFlags() error{
	switch f.resource{
	case "item":
		return f.handleItemCommand()
	case "tag":
		return f.handleTagCommands()
	}
	return nil
}



