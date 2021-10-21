package helpers

import (
	"reflect"
	"strings"
)

func GetMapKeys(object interface{}) []string {
	reflectValue := reflect.ValueOf(object)
	if reflectValue.Kind() != reflect.Map {
		panic("value must be a map")
	}

	reflectType := reflect.TypeOf(object)
	if reflectType.Key().Kind() != reflect.String {
		panic("key must be a string")
	}

	keys := []string{}
	for _, key := range reflectValue.MapKeys() {
		keys = append(keys, key.String())
	}

	return keys
}

func GetMapValues(object interface{}) []interface{} {
	reflectValue := reflect.ValueOf(object)
	if reflectValue.Kind() != reflect.Map {
		panic("value must be a map")
	}

	iter := reflectValue.MapRange()
	values := []interface{}{}
	for iter.Next() {
		values = append(values, iter.Value().Interface())
	}

	return values
}

func GetStructFields(value interface{}, exclude []interface{}) []interface{} {
	reflectionType := reflect.TypeOf(value)

	fields := []interface{}{}
	for i := 0; i < reflectionType.NumField(); i++ {
		field := strings.Split(reflectionType.Field(i).Tag.Get("json"), ",")[0]
		if !Contains(exclude, field) {
			fields = append(fields, field)
		}
	}

	return fields
}

func GenerateUserReturnColumns(excludedColumns []string) []string {
	allColumns := []string{
		"id",
		"average_rating",
		"created_at",
		"email",
		"firstname",
		"image",
		`CASE 
  			WHEN otp_secret_key = '' THEN CAST ('false' AS BOOLEAN)
  			ELSE CAST('true' AS BOOLEAN)
			END AS is_2fa_enabled`,
		`CASE 
  			WHEN email_verified_at IS NULL THEN CAST ('false' AS BOOLEAN)
  			ELSE CAST('true' AS BOOLEAN)
			END AS is_email_verified`,
		`CASE 
  			WHEN phone_verified_at IS NULL THEN CAST ('false' AS BOOLEAN)
  			ELSE CAST('true' AS BOOLEAN)
			END AS is_phone_verified`,
		"lastname",
		"password",
		"phone_no",
		"reviews_count",
		"trips_count",
	}

	includedColumns := []string{}
	for _, returnColumn := range allColumns {
		if !ContainsSuffix(excludedColumns, returnColumn) {
			includedColumns = append(includedColumns, returnColumn)
		}
	}
	return includedColumns
}
