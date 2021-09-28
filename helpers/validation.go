package helpers

// Overiding the gin framework's default validator which implements the StructValidator interface.

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	lowercaseRegex        = regexp.MustCompile(`[a-z]+`)
	nameRegex             = regexp.MustCompile(`^[a-zA-z ]+$`)
	numberRegex           = regexp.MustCompile(`\d+`)
	specialCharacterRegex = regexp.MustCompile(`\W+`)
	uppercaseRegex        = regexp.MustCompile(`[A-Z]+`)
)

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) != reflect.Struct {
		return nil
	}

	v.lazyinit()
	return v.validate.Struct(obj)
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		v.validate.RegisterTagNameFunc(jsonTagName)

		err := v.validate.RegisterValidation("name", validateName)
		ExitIfError(err)

		err = v.validate.RegisterValidation("password", validatePassword)
		ExitIfError(err)
	})
}

func jsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return specialCharacterRegex.MatchString(value) &&
		lowercaseRegex.MatchString(value) &&
		uppercaseRegex.MatchString(value) &&
		numberRegex.MatchString(value)
}
