package command

type Command interface {
	Name() string
	Exec(args []string) error
	parseFlags(args []string) error
}


type BaseCommand struct {
	name        string
	flags *Flags
}

func (cmd *BaseCommand) Name() string        { return cmd.name }