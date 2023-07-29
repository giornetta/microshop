package pg

import (
	"context"
	"fmt"

	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewIdentityRepository(pool *pgxpool.Pool) identity.IdentityRepository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) ExistsByEmail(email string, ctx context.Context) (bool, error) {
	if err := r.pool.QueryRow(ctx, "SELECT 1 FROM identities WHERE email = $1", email).Scan(); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, &errors.ErrInternal{Err: err}
	}

	return true, nil
}

func (r *repository) FindByEmail(email string, ctx context.Context) (*identity.Identity, error) {
	var ident identity.Identity

	if err := r.pool.QueryRow(
		ctx,
		"SELECT identity_id, email, roles, password_hash FROM identities WHERE email = $1",
		email,
	).Scan(&ident.Id, &ident.Email, &ident.Roles, &ident.PasswordHash); err != nil {
		if err == pgx.ErrNoRows {
			return nil, &identity.ErrNotFound{Email: email}
		}
		fmt.Println(err)
		return nil, &errors.ErrInternal{Err: err}
	}

	return &ident, nil
}

func (r *repository) FindById(id identity.IdentityId, ctx context.Context) (*identity.Identity, error) {
	var ident identity.Identity

	if err := r.pool.QueryRow(
		ctx,
		"SELECT identity_id, email, roles, password_hash FROM identities WHERE identity_id = $1",
		id,
	).Scan(&ident.Id, &ident.Email, &ident.Roles, &ident.PasswordHash); err != nil {
		if err == pgx.ErrNoRows {
			return nil, &identity.ErrNotFound{IdentityId: id}
		}

		return nil, &errors.ErrInternal{Err: err}
	}

	return &ident, nil
}

func (r *repository) Store(identity *identity.Identity, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		`INSERT INTO identities(identity_id, email, roles, password_hash) 
		 VALUES($1, $2, $3, $4);`,
		identity.Id, identity.Email, identity.Roles, identity.PasswordHash,
	); err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(identity *identity.Identity, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		`UPDATE identities 
		 SET email = $1, roles = $2, password_hash = $3
		 WHERE identity_id = $4`,
		identity.Email, identity.Roles, identity.PasswordHash, identity.Id.String(),
	); err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(id identity.IdentityId, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		`DELETE FROM identities WHERE identity_id = $1`,
		id.String(),
	); err != nil {
		return err
	}

	return nil
}
