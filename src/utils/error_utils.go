package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func SplitError(validationError error) []string {
	// Create New Error Array Variable
	var errors []string

	// Check Error Type
	if strings.Contains(validationError.Error(), "Error:Field") {
		// Loop Through Error
		for _, err := range validationError.(validator.ValidationErrors) {
			// Create New Error String
			errorString := err.Field() + " " + err.Tag()
			// Append Error to Array
			errors = append(errors, errorString)
		}
	} else {
		errors = append(errors, validationError.Error())
	}

	// Return errors
	return errors
}
