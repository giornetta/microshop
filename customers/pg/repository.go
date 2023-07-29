package pg

import (
	"context"
	"database/sql"

	"github.com/giornetta/microshop/customers"
	"github.com/giornetta/microshop/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type customerModel struct {
	Id string

	FirstName string
	LastName  string

	Country sql.NullString
	City    sql.NullString
	ZipCode sql.NullString
	Street  sql.NullString
}

func (m customerModel) ToCustomer() *customers.Customer {
	return &customers.Customer{
		Id:        customers.CustomerId(m.Id),
		FirstName: m.FirstName,
		LastName:  m.LastName,
		ShippingAddress: &customers.ShippingAddress{
			Country: m.Country.String,
			City:    m.City.String,
			ZipCode: m.ZipCode.String,
			Street:  m.Street.String,
		},
	}
}

type repository struct {
	pool *pgxpool.Pool
}

func NewCustomerRepository(pool *pgxpool.Pool) customers.CustomerRepository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Store(c *customers.Customer, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		`INSERT INTO
		customers(customer_id, first_name, last_name)
		VALUES($1, $2, $3);`,
		c.Id, c.FirstName, c.LastName,
	); err != nil {
		return err
	}

	return nil
}

func (r *repository) FindById(id customers.CustomerId, ctx context.Context) (*customers.Customer, error) {
	var c customerModel

	if err := r.pool.QueryRow(
		ctx,
		`SELECT customer_id, first_name, last_name, shipping_country, shipping_city, shipping_zipcode, shipping_street
		FROM customers WHERE customer_id = $1`,
		id,
	).Scan(
		&c.Id, &c.FirstName, &c.LastName,
		&c.Country, &c.City, &c.ZipCode, &c.Street,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, &customers.ErrNotFound{CustomerId: id}
		}

		return nil, &errors.ErrInternal{Err: err}
	}

	return c.ToCustomer(), nil
}

func (r *repository) UpdateShippingAddress(id customers.CustomerId, addr *customers.ShippingAddress, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		`UPDATE customers
		SET shipping_country = $1, shipping_city = $2, shipping_zipcode = $3, shipping_street = $4
		WHERE customer_id = $5;`,
		addr.Country, addr.City, addr.ZipCode, addr.Street, id,
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
