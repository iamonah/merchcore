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
			for _, verr := range validationErrors {
				switch verr.Tag() {
				case "required":
					fe.AddFieldError(strings.ToLower(verr.Field()), errors.New("field is required"))
				}
			}
		}
	}
	return fe.ToError()
}
