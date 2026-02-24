package validator

import (
	"fmt"
	"net/mail"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type Errors []ValidationError

func (e Errors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, ", ")
}

func (e Errors) HasErrors() bool {
	return len(e) > 0
}

type Validator struct {
	errors Errors
}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{Field: field, Message: "is required"})
	}
	return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(strings.TrimSpace(value)) < min {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at least %d characters", min),
		})
	}
	return v
}

func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(strings.TrimSpace(value)) > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at most %d characters", max),
		})
	}
	return v
}

func (v *Validator) Email(field, value string) *Validator {
	if _, err := mail.ParseAddress(value); err != nil {
		v.errors = append(v.errors, ValidationError{Field: field, Message: "is not a valid email"})
	}
	return v
}

func (v *Validator) OneOf(field, value string, allowed ...string) *Validator {
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")),
	})
	return v
}

func (v *Validator) Validate() error {
	if v.errors.HasErrors() {
		return v.errors
	}
	return nil
}