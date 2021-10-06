package helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GenerateErrorMessages(errors validator.ValidationErrors) map[string]string {
	messages := make(map[string]string)

	for _, err := range errors {
		field := err.Field()
		switch err.ActualTag() {
		case "name":
			messages[field] = fmt.Sprintf("%v should contain only letters and spaces", strings.Title(field))
		case "email":
			messages[field] = fmt.Sprintf("%v is not a valid email address", strings.Title(field))
		case "gt":
			messages[field] = fmt.Sprintf("%v should be greater than %v", strings.Title(field), err.Param())
		case "lte":
			messages[field] = fmt.Sprintf("%v should not be greater than %v", strings.Title(field), err.Param())
		case "max":
			messages[field] = fmt.Sprintf("%v should not be greater than %v characters", strings.Title(field), err.Param())
		case "min":
			messages[field] = fmt.Sprintf("%v should not be less than %v characters", strings.Title(field), err.Param())
		case "oneof":
			messages[field] = fmt.Sprintf("%v should be in these category %v", strings.Title(field), err.Param())
		case "password":
			messages[field] = fmt.Sprintf("%v should be a mix of uppercase, lowercase, numeric and special characters", strings.Title(field))
		case "required":
			if field == "token" {
				// Make things easier on the frontend
				messages["message"] = fmt.Sprintf("%v is required", strings.Title(field))
			} else {
				messages[field] = fmt.Sprintf("%v is required", strings.Title(field))
			}
		default:
			messages[field] = fmt.Sprintf("%v is invalid", strings.Title(field))
		}
	}

	return messages
}
