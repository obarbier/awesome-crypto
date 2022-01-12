package common

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	// register function to get tag name from json tags.
	// https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go#L36
	//v.RegisterTagNameFunc(func(fld reflect.StructField) string {
	//	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	//	if name == "-" {
	//		return ""
	//	}
	//	return name
	//})

	err := v.RegisterValidation("not-blank", validators.NotBlank)
	if err != nil {
		// TODO: handle Error
	}
	return v
}
