package cmdflag

import (
	"fmt"

	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

func (f *Flags) handleTagCommands() error {

	switch f.command {
	case "view":
		return f.TagGet()
	case "delete":
		return f.TagDelete()
	default:
		return fmt.Errorf("unknown action (%s) for Tag", f.command)
	}
}

func (f *Flags) TagGet() error {
	// tags, err := tcmd.repo.GetTags()
	// if err != nil {
	// 	return err
	// }
	// // response.Respond(f.flags.Format, tags)

	return nil
}
func (f *Flags) TagDelete() error {

	// err = tcmd.repo.DeleteTag(tcmd.flags.Name)

	return validation.ValidateFlagsDefinedStr([]string{"-n"}, f.Name)
}
