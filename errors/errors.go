package errors

import "net/http"

type WithStatusCode interface {
	error
	StatusCode() int
}

type ErrInternal struct {
	Err error
}

func (e *ErrInternal) Error() string {
	return "Internal Server Error"
}

func (e *ErrInternal) Cause() error {
	return e.Err
}

func (e *ErrInternal) StatusCode() int {
	return http.StatusInternalServerError
}

type ErrBadRequest struct{}

func (e *ErrBadRequest) Error() string {
	return "Bad Request"
}

func (e *ErrBadRequest) StatusCode() int {
	return http.StatusBadRequest
}
