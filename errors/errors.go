package errors

import (
	"net/http"
)

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

type ErrBadRequest struct {
	Err error
}

func (e *ErrBadRequest) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return "Bad Request"
}

func (e *ErrBadRequest) Cause() error {
	return e.Err
}

func (e *ErrBadRequest) StatusCode() int {
	return http.StatusBadRequest
}

type ErrUnauthorized struct {
	Err error
}

func (e *ErrUnauthorized) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return "Unauthorized"
}

func (e *ErrUnauthorized) Cause() error {
	return e.Err
}

func (e *ErrUnauthorized) StatusCode() int {
	return http.StatusUnauthorized
}
