package validation

import (
	"fmt"
	"strings"
)

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

func ValidateResourceAction(cmd map[string][]string, resource string, action string) error {
	_, ok := cmd[resource]

	if !ok {
		return &ValidationError{
			Field:   "resource",
			Message: fmt.Sprintf("Unknown resource %s", resource),
		}
	}

	for _, act := range cmd[resource] {
		if action == act {
			return nil
		}
	}
	return &ValidationError{
		Field:   "action",
		Message: fmt.Sprintf("Unknown action %s for resource %s", action, resource),
	}
}

func ValidateFlagCount(provided int, required int) error {
	if provided < required {
		return &ValidationError{
			Field:   "arguments",
			Message: fmt.Sprintf("insufficient number of arguments: expected at least %d, got %d", required, provided),
		}
	}
	return nil
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
