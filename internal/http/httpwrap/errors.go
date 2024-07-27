package httpwrap

import (
	"errors"
	"fmt"
	"net/http"
)

type statusWrapper struct {
	err        error
	statusCode int
}

func (w statusWrapper) Error() string {
	return w.err.Error()
}

func (w statusWrapper) Unwrap() error {
	return w.err
}

func StatusCodeError(err error, code int) error {
	// Don't discard any status already here.
	if errors.As(err, &statusWrapper{}) {
		return err
	}
	return statusWrapper{
		err:        err,
		statusCode: code,
	}
}

func ErrBadRequestf(msg string, args ...any) error {
	return StatusCodeError(fmt.Errorf(msg, args...), http.StatusBadRequest)
}

func StatusCodeForError(err error) int {
	var w statusWrapper
	if errors.As(err, &w) {
		return w.statusCode
	}

	return http.StatusInternalServerError
}
