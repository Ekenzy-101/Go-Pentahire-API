package helpers

func Contains(slice []interface{}, value interface{}) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}
