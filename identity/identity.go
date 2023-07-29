package identity

import (
	"context"

	"github.com/giornetta/microshop/auth"
)

type IdentityId string

func (id IdentityId) String() string {
	return string(id)
}

type Identity struct {
	Id           IdentityId
	Email        string
	PasswordHash []byte
	Roles        []auth.Role
}

type IdentityQuerier interface {
	FindById(id IdentityId, ctx context.Context) (*Identity, error)
	FindByEmail(email string, ctx context.Context) (*Identity, error)
	ExistsByEmail(email string, ctx context.Context) (bool, error)
}

type IdentityStorer interface {
	Store(identity *Identity, ctx context.Context) error
	Update(identity *Identity, ctx context.Context) error
	Delete(id IdentityId, ctx context.Context) error
}

type IdentityRepository interface {
	IdentityQuerier
	IdentityStorer
}

type Service interface {
	SignUp(req *SignUpRequest, ctx context.Context) (*SignUpResponse, error)
	SignIn(req *SignInRequest, ctx context.Context) (*SignInResponse, error)
	UpdateRoles(req *UpdateRolesRequest, ctx context.Context) error
}

type SignUpRequest struct {
	Email    string
	Password string
}

type SignUpResponse struct {
	IdentityId IdentityId
	Token      string
}

type SignInRequest struct {
	Email    string
	Password string
}

type SignInResponse struct {
	Token string
}

type UpdateRolesRequest struct {
	Id    IdentityId
	Roles []auth.Role
}
