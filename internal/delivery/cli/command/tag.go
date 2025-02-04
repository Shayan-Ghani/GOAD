package command

import (
	"flag"
	"fmt"
	"github.com/Shayan-Ghani/GOAD/internal/repository"
	"github.com/Shayan-Ghani/GOAD/pkg/response"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

type TagCommand struct {
	BaseCommand
	repo repository.Repository
}

func NewTagCommand(repo repository.Repository, action string) *TagCommand {
	return &TagCommand{
		BaseCommand: BaseCommand{
			name:        action,
			flags:       &Flags{},
		},
		repo: repo,
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
	fs.StringVar(&tcmd.flags.Format, "format", "table", fmt.Sprintf("%s output format", resource))

	err := fs.Parse(args[2:])
	if args[1] == "--help"{
		fs.PrintDefaults()
	}

	return err
}

func (tcmd *TagCommand) handleView() error {
	tags, err := tcmd.repo.GetTags()
	if err != nil {
		return err
	}
	response.Respond(tcmd.flags.Format, tags)

	return err
}
func (tcmd *TagCommand) handleDelete() error {
	var err error
	if err = validation.ValidateFlagsDefinedStr([]string{"-n"}, tcmd.flags.Name); err != nil {
		return fmt.Errorf("%w", err)
	}
	err = tcmd.repo.DeleteTag(tcmd.flags.Name)

	return err
}
