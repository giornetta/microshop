package auth

import (
	"context"
	"errors"
)

var (
	ErrExpired       error = errors.New("token is expired")
	ErrInvalidIssuer error = errors.New("token issuer is invalid")
)

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	CustomerRole Role = "Customer"
	AdminRole    Role = "Admin"
	CourierRole  Role = "Courier"
)

type contextKey string

const ContextKey contextKey = "AUTH_TOKEN"

func FromContext(ctx context.Context) (Token, error) {
	val := ctx.Value(ContextKey)
	switch t := val.(type) {
	case *jwtToken:
		return t, nil
	default:
		return nil, errors.New("no token")
	}
}

type Token interface {
	Subject() string
	Roles() []Role
	IsAdmin() bool
}

type Issuer interface {
	Issue(subject string, roles []Role) (string, error)
}

type Verifier interface {
	Verify(signed string) (Token, error)
}
