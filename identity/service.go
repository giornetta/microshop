package identity

import (
	"context"
	"fmt"

	"github.com/giornetta/microshop/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	issuer     auth.Issuer
	repository IdentityRepository
}

func NewService(issuer auth.Issuer, repository IdentityRepository) Service {
	return &service{
		issuer:     issuer,
		repository: repository,
	}
}

func (s *service) SignUp(req *SignUpRequest, ctx context.Context) (*SignUpResponse, error) {
	// TODO Validate Inputs

	exists, err := s.repository.ExistsByEmail(req.Email, ctx)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := &Identity{
		Id:           IdentityId(uuid.New().String()),
		Email:        req.Email,
		PasswordHash: hash,
		Roles:        []auth.Role{auth.CustomerRole},
	}

	if err := s.repository.Store(id, ctx); err != nil {
		return nil, err
	}

	token, err := s.issuer.Issue(id.Id.String(), id.Roles)
	if err != nil {
		return nil, err
	}

	return &SignUpResponse{
		IdentityId: id.Id,
		Token:      token,
	}, nil
}

func (s *service) SignIn(req *SignInRequest, ctx context.Context) (*SignInResponse, error) {
	// TODO Validate

	ident, err := s.repository.FindByEmail(req.Email, ctx)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(ident.PasswordHash, []byte(req.Password)); err != nil {
		return nil, err
	}

	tok, err := s.issuer.Issue(ident.Id.String(), ident.Roles)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{
		Token: tok,
	}, nil
}

func (s *service) UpdateRoles(req *UpdateRolesRequest, ctx context.Context) error {
	// TODO Validate

	ident, err := s.repository.FindById(req.Id, ctx)
	if err != nil {
		return err
	}

	ident.Roles = req.Roles

	if err := s.repository.Update(ident, ctx); err != nil {
		return err
	}

	return nil
}
