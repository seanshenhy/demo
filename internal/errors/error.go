package errors

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
)

const UNAUTHORIZED = 401
const FORBIDDEN = 403
const NOT_FOUND = 404
const PARAM_INVALID = 422

type HTTPError struct {
	Errors map[string][]string `json:"errors"`
	Code   int                 `json:"-"`
}

func NewHttpError(code int, field string, detail string) *HTTPError {
	return &HTTPError{
		Code: code,
		Errors: map[string][]string{
			field: {detail},
		},
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError code: %d message: %v", e.Code, e.Errors)
}

func FromError(err error) *HTTPError {
	if err == nil {
		return nil
	}
	if se := new(errors.Error); errors.As(err, &se) {
		if se.Reason == "CODEC" {
			return NewHttpError(int(se.Code), "message", se.Message)
		}
		return NewHttpError(int(se.Code), se.Reason, se.Message)
	}
	return NewHttpError(500, "internal", "error")
}
