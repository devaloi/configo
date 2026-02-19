package configo

import (
	"fmt"
	"strings"
)

// KeyNotFoundError indicates a requested key does not exist.
type KeyNotFoundError struct {
	Key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found: %s", e.Key)
}

// TypeMismatchError indicates a value cannot be converted to the requested type.
type TypeMismatchError struct {
	Key      string
	Expected string
	Actual   any
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch for key %q: expected %s, got %T", e.Key, e.Expected, e.Actual)
}

// FieldError holds a validation error for a single field.
type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationError collects multiple field-level errors.
type ValidationError struct {
	Errors []FieldError
}

func (e *ValidationError) Error() string {
	msgs := make([]string, len(e.Errors))
	for i, fe := range e.Errors {
		msgs[i] = fe.Error()
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, "; "))
}
