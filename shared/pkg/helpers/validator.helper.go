package helpers

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func validateStruct(param any) []ValidationError {
	var errs []ValidationError
	validate := validator.New()
	err := validate.Struct(param)

	if err != nil {
		v := reflect.ValueOf(param)
		t := v.Type()

		for _, err := range err.(validator.ValidationErrors) {
			field, _ := t.FieldByName(err.StructField())
			formName := field.Tag.Get("form")
			if formName == "" {
				formName = err.StructField()
			}

			element := ValidationError{
				Field: formName,
				Tag:   err.Tag(),
			}
			errs = append(errs, element)
		}
	}
	return errs
}

func ValidateInput(input interface{}) *[]ValidationError {
	if errs := validateStruct(input); len(errs) > 0 {
		return &errs
	}

	return nil
}
