package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
)

type AppErr struct {
	Code     int         `json:"code"`
	Message  string      `json:"message,omitempty"`
	Fields   FieldErrors `json:"fields,omitempty"`
	FuncName string      `json:"-"`
	FileName string      `json:"-"`
}

func (e *AppErr) Error() string {
	byte, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(byte)
}

func New(code ErrCode, err error) *AppErr {
	pc, filename, line, _ := runtime.Caller(1)

	var fields *FieldErrors
	if errors.As(err, &fields) && len(*fields) > 0 {
		return &AppErr{
			Code:     HTTPStatus[code],
			Message:  "Validation failed",
			Fields:   *fields,
			FuncName: runtime.FuncForPC(pc).Name(),
			FileName: fmt.Sprintf("%s:%d", filename, line),
		}
	}

	return &AppErr{
		Code:     HTTPStatus[code],
		Message:  err.Error(),
		FuncName: runtime.FuncForPC(pc).Name(),
		FileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

func Newf(code ErrCode, format string, v ...any) *AppErr {
	pc, filename, line, _ := runtime.Caller(1)

	return &AppErr{
		Code:     HTTPStatus[code],
		Message:  fmt.Errorf(format, v...).Error(),
		FuncName: runtime.FuncForPC(pc).Name(),
		FileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

type DomainError struct {
	Msg  error
	Code ErrCode
}

func (e *DomainError) Error() string {
	return e.Msg.Error()
}

func (e *DomainError) Unwrap() error {
	return e.Msg
}

func NewDomainError(code ErrCode, msg error) error {
	err := fmt.Errorf("%w", msg)
	return &DomainError{Code: code, Msg: err}
}

func IsDomainError(err error) (*DomainError, bool) {
	var dError *DomainError
	if errors.As(err, &dError) {
		return dError, true
	}
	return nil, false
}

type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

type FieldErrors []FieldError

func NewFieldErrors() *FieldErrors {
	return &FieldErrors{}
}

func (fe *FieldErrors) AddFieldError(field string, err error) {
	*fe = append(*fe, FieldError{
		Field: field,
		Err:   err.Error(),
	})
}

func (fe *FieldErrors) ToError() error {
	if len(*fe) != 0 {
		return fe
	}
	return nil
}

func (fe *FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// func (fe FieldErrors) MarshalJSON() ([]byte, error) {
// 	return json.Marshal([]FieldError(fe))
// }
