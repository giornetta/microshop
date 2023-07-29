package identity

import (
	"fmt"
	"net/http"
)

type ErrNotFound struct {
	IdentityId IdentityId
	Email      string
}

func (err *ErrNotFound) Error() string {
	if err.IdentityId != "" {
		return fmt.Sprintf("user with id=%s was not found", err.IdentityId.String())
	}

	if err.Email != "" {
		return fmt.Sprintf("user with email=%s was not found", err.Email)
	}

	return "user was not found"
}

func (err *ErrNotFound) StatusCode() int {
	return http.StatusNotFound
}

type ErrAlreadyExists struct {
	Email string
}

func (err *ErrAlreadyExists) Error() string {
	return fmt.Sprintf("user with email=%s already exists", err.Email)
}

func (err *ErrAlreadyExists) StatusCode() int {
	return http.StatusBadRequest
}
