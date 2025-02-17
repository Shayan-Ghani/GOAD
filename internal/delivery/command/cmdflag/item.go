package cmdflag

import (
	"fmt"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

func (f *Flags) handleItemCommand() error {

	switch f.command {
	case "add":
		return f.ItemAdd()
	case "view":
		return f.ItemView()
	case "delete":
		return f.ItemDelete()
	case "update":
		return f.ItemUpdate()
	case "done":
		return f.ItemDone()
	default:
		return nil
	}
}

func (f *Flags) ItemAdd() error {
	return validation.ValidateFlagsDefinedStr([]string{"-n", "-d"}, f.Name, f.Description)
}

func (f *Flags) ItemView() error {
	return nil
}

func (f *Flags) ItemUpdate() error {

	if err := validation.ValidateFlagsDefinedStr([]string{"-i"}, f.ID); err != nil {
		return err
	}

	if f.Description == "" &&  f.Name == "" && f.Tags == "" {
		return fmt.Errorf("description, name or tags must be defined")
	}

	return nil
}

func (f *Flags) ItemDone() error {
	return validation.ValidateFlagsDefinedStr([]string{"-i"}, f.ID)
}

func (f *Flags) ItemDelete() error {
	return validation.ValidateFlagsDefinedStr([]string{"-i"}, f.ID)
}