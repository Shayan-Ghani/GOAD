package command

import (
	cmdflag "github.com/Shayan-Ghani/GOAD/internal/delivery/command/cmdflag"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

type Command struct {
	resource string
	command  string
	flags    *cmdflag.Flags
}

func NewCommand(args []string) (*Command, error) {
	if len(args) < 1 {
		return nil, validation.Help{}
	}

	if Ishelp(args[0]) {
		return nil, validation.Help{}
	}

	if err := validation.ValidateCommand(args); err != nil {
		return nil, err
	}

	resource := args[0]
	command := args[1]

	f := cmdflag.New(resource, command)

	if err := f.Parse(args); err != nil {
		return nil, err
	}

	if err := f.CheckCommandFlags(); err != nil {
		return nil, err
	}

	return &Command{
		resource: args[0],
		command:  args[1],
		flags:    f,
	}, nil
}

func (c Command) GetFlags() *cmdflag.Flags {
	return c.flags
}

func Ishelp(arg string) bool {
	if arg == "--help" || arg == "-h" {
		return true
	}
	return false
}
