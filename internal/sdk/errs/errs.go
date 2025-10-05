package errs

import (
	"encoding/json"
	"fmt"
)

type AppErr struct {
	Code    int    `json:"code"`
	Err     string `json:"error"`
	Details string `json:"details,omitempty"`
}

func (appErr *AppErr) String() string {
	return fmt.Sprintf("error: %s, details: %s", appErr.Err, appErr.Details)
}

func (appErr *AppErr) Error() string {
	d, err := json.Marshal(appErr)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"message"`
}

type FieldErrors []FieldError

func NewFieldErrors() *FieldErrors {
	return &FieldErrors{}
}

func (fe *FieldErrors) AddFieldError(field string, err error) {
	newerror := FieldError{
		Field: field,
		Err:   err.Error(),
	}
	*fe = append(*fe, newerror)
}

func (fe *FieldErrors) ToError() *AppErr {
	if len(*fe) != 0 {
		return NewAppErr(InvalidArgument, fe)
	}
	return nil
}

func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

func NewAppErr(code ErrCode, err error) *AppErr {
	name := CodeNames[code]
	statusCode := HTTPStatus[code]

	return &AppErr{
		Code:    statusCode,
		Err:     name,
		Details: err.Error(),
	}
}
