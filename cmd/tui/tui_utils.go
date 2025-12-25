package tui

import "fmt"

// ValidationError satisfies the Error interface and represents a type of
// validation error for checking user inputs
type ValidationError struct {
	Field   string
	Message string
}

// Error returns the string representation of the error
func (va *ValidationError) Error() string {
	return fmt.Sprintf("Error with field %s: %s", va.Field, va.Message)
}
