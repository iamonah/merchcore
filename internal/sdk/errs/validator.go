package errs

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func NewValidate(data any) error {
	err := Validate.Struct(data)
	fe := NewFieldErrors()

	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range validationErrors {
				switch err.Tag() {
				case "required":
					fe.AddFieldError(strings.ToLower(err.Field()), errors.New("field is required"))
				}
			}
		}
	}

	if len(*fe) != 0 {
		return fe
	}
	return nil
}
