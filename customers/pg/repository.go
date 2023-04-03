package pg

import (
	"context"

	"github.com/giornetta/microshop/customers"
	"github.com/giornetta/microshop/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) customers.CustomerRepository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Store(c *customers.Customer, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		`INSERT INTO
		customers(customer_id, first_name, last_name, email, shipping_country, shipping_city, shipping_zipcode, shipping_street)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,
		c.Id, c.FirstName, c.LastName, c.Email,
		c.ShippingAddress.Country, c.ShippingAddress.City, c.ShippingAddress.ZipCode, c.ShippingAddress.ZipCode,
	); err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByEmail(email string, ctx context.Context) (*customers.Customer, error) {
	var c customers.Customer

	if err := r.pool.QueryRow(
		ctx,
		`SELECT customer_id, first_name, last_name, email, shipping_country, shipping_city, shipping_zipcode, shipping_street
		FROM customers WHERE email = $1`,
		email,
	).Scan(
		&c.Id, &c.FirstName, &c.LastName, &c.Email,
		&c.ShippingAddress.Country, &c.ShippingAddress.City, &c.ShippingAddress.ZipCode, &c.ShippingAddress.Street,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, &customers.ErrNotFound{Email: email}
		}

		return nil, &errors.ErrInternal{Err: err}
	}

	return &c, nil
}

func (r *repository) FindById(id customers.CustomerId, ctx context.Context) (*customers.Customer, error) {
	var c customers.Customer

	if err := r.pool.QueryRow(
		ctx,
		`SELECT customer_id, first_name, last_name, email, shipping_country, shipping_city, shipping_zipcode, shipping_street
		FROM customers WHERE customer_id = $1`,
		id,
	).Scan(
		&c.Id, &c.FirstName, &c.LastName, &c.Email,
		&c.ShippingAddress.Country, &c.ShippingAddress.City, &c.ShippingAddress.ZipCode, &c.ShippingAddress.Street,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, &customers.ErrNotFound{CustomerId: id}
		}

		return nil, &errors.ErrInternal{Err: err}
	}

	return &c, nil
}

func (r *repository) Update(c *customers.Customer, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		`UPDATE customers
		SET shipping_country = $1, shipping_city = $2, shipping_zipcode = $3, shipping_street = $4
		WHERE customer_id = $5;`,
		c.ShippingAddress.Country, c.ShippingAddress.City, c.ShippingAddress.ZipCode, c.ShippingAddress.Street, c.Id,
	); err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(id customers.CustomerId, ctx context.Context) error {
	if _, err := r.pool.Exec(ctx, "DELETE FROM customers WHERE customer_id = $1;", id); err != nil {
		return err
	}

	return nil
}
