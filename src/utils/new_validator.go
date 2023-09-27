package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validate *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()
	//validate.RegisterValidation("required_unless_null", requiredUnlessNull)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return &CustomValidator{Validate: validate}
}

func (v *CustomValidator) ValidateStruct(value interface{}) (errs []string) {
	errValidation := v.Validate.Struct(value)
	if errValidation != nil {
		for _, err := range errValidation.(validator.ValidationErrors) {
			// Create error messages based on errors returned ...
			newErrFormat := fmt.Sprintf(
				"'%s': failed '%s' tag check (value '%s' is not valid)",
				err.Field(), err.Tag(), err.Value(),
			)
			errs = append(errs, newErrFormat)
		}
		return
	}
	return nil
}

func requiredUnlessNull(fl validator.FieldLevel) bool {
	field := fl.Field().Interface()
	param := fl.Param()
	specifiedFieldValue := fl.Parent().FieldByName(param).Interface()

	return specifiedFieldValue == nil || (reflect.TypeOf(specifiedFieldValue).Kind() == reflect.Ptr && reflect.ValueOf(specifiedFieldValue).IsNil()) || (field != nil && field != "")
}
