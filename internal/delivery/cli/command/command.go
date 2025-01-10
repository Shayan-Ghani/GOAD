package command

type Command interface {
	Name() string
	Exec(args []string) error
}


type BaseCommand struct {
	name        string
	description string
	usage       string
	flags *Flags
}

func (cmd *BaseCommand) Name() string        { return cmd.name }