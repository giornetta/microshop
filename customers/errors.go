package customers

import (
	"fmt"
	"net/http"
)

type ErrNotFound struct {
	CustomerId CustomerId
	Email      string
}

func (err *ErrNotFound) Error() string {
	if err.CustomerId != "" {
		return fmt.Sprintf("customer with id=%s was not found", err.CustomerId.String())
	}

	if err.Email != "" {
		return fmt.Sprintf("customer with email=%s was not found", err.Email)
	}

	return "customer was not found"
}

func (err *ErrNotFound) StatusCode() int {
	return http.StatusNotFound
}

type ErrAlreadyExists struct {
	Id CustomerId
}

func (err *ErrAlreadyExists) Error() string {
	return fmt.Sprintf("customer with id=%s already exists", err.Id)
}

func (err *ErrAlreadyExists) StatusCode() int {
	return http.StatusBadRequest
}
