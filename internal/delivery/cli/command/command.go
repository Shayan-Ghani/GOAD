package command

type Command interface {
	Name() string
	Description() string
	Usage() string
	Exec(args []string) error
}


type BaseCommand struct {
	name        string
	description string
	usage       string
	flags *Flags
}

func (cmd *BaseCommand) Name() string        { return cmd.name }
func (cmd *BaseCommand) Description() string { return cmd.description }
func (cmd *BaseCommand) Usage() string       { return cmd.usage }
