package validation

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
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