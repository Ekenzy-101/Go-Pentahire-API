package helpers

import "strings"

func Contains(slice []interface{}, value interface{}) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

// Checks if value ends with any item of the slice
func ContainsSuffix(slice []string, value string) bool {
	for _, item := range slice {
		if strings.HasSuffix(value, item) {
			return true
		}
	}
	return false
}
