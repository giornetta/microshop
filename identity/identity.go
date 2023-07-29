package identity

import (
	"context"
	"regexp"
	"strings"

	"github.com/giornetta/microshop/auth"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

func (r *SignUpRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)

	return validation.ValidateStruct(r,
		validation.Field(&r.Email,
			validation.Required,
			is.EmailFormat,
		),
		validation.Field(&r.Password,
			validation.Required,
			validation.Length(8, 24),
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9!#@?]+")),
		),
	)
}

type SignUpResponse struct {
	IdentityId IdentityId
	Token      string
}

type SignInRequest struct {
	Email    string
	Password string
}

func (r *SignInRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)

	return validation.ValidateStruct(r,
		validation.Field(&r.Email,
			validation.Required,
			is.EmailFormat,
		),
		validation.Field(&r.Password,
			validation.Required,
			validation.Length(8, 24),
		),
	)
}

type SignInResponse struct {
	Token string
}

type UpdateRolesRequest struct {
	Id    IdentityId
	Roles []auth.Role
}

func (r *UpdateRolesRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Id,
			validation.Required,
		),
		validation.Field(&r.Roles,
			validation.Each(
				validation.In(auth.CustomerRole, auth.AdminRole, auth.CourierRole),
			),
		),
	)
}
