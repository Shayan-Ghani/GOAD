package validation

import (
	"fmt"
	"strings"
)


type Help struct{
	Message string
}

func (h Help) Error() string {
    return "help requested"
}

type ValidationError struct {
	Field   string
	Message string
}

func New(field string, msg string) error {
	return &ValidationError{
		Field:   field,
		Message: msg,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateTagName(tag string) error {
	if strings.TrimSpace(tag) == "" {
		return &ValidationError{
			Field:   "tag",
			Message: "tag name cannot be empty",
		}
	}

	if strings.ContainsAny(tag, "!@#$%^&*()") {
		return &ValidationError{
			Field:   "tag",
			Message: "tag name contains invalid characters",
		}
	}

	return nil
}

func ValidateTagNames(tags []string) error {
	for _, tag := range tags {
		if err := ValidateTagName(tag); err != nil {
			return err
		}
	}
	return nil
}

func ValidateArgCount(args []string) error {

	l := len(args)
	ValidCount := 2

	if l < ValidCount {
		return &ValidationError{
			Field:   "arguments",
			Message: fmt.Sprintf("insufficient number of arguments: expected at least %d, got %d", ValidCount, l),
		}
	}

	ValidCount = 4

	if command := args[1]; command == "update" &&  l < ValidCount {
		return &ValidationError{
			Field:   "arguments",
			Message: fmt.Sprintf("insufficient number of arguments for update: expected at least %d, got %d", ValidCount, l),
		}
	}
	return nil
}


func ValidateCommand(args []string) error {
	if err := ValidateArgCount(args); err != nil {
		return err
	}

	resource := args[0]
	command := args[1]

	validCommands := map[string][]string{
		"item": {"add", "view", "delete", "update", "done", "--help"},
		"tag":  {"view", "delete", "--help"},
	}

	_, ok := validCommands[resource]

	if !ok {
		return &ValidationError{
			Field:   "resource",
			Message: fmt.Sprintf("Unknown resource %s", resource),
		}
	}

	for _, act := range validCommands[resource] {
		if command == act {
			return nil
		}
	}
	return &ValidationError{
		Field:   "command",
		Message: fmt.Sprintf("Unknown command %s for resource %s", command, resource),
	}
}

func ValidateFlagsDefinedStr(argNames []string, flags ...string) error {
	for i, f := range flags {
		if f == "" {
			return &ValidationError{
				Field:   "argument",
				Message: fmt.Sprintf("%s is required!", argNames[i]),
			}
		}
	}
	return nil
}
