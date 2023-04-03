package pg

import (
	"context"

	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/products"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) products.ProductRepository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) FindById(id products.ProductId, ctx context.Context) (*products.Product, error) {
	var product products.Product

	err := r.pool.QueryRow(ctx, "SELECT * FROM products WHERE product_id = $1", id).Scan(
		&product.Id, &product.Name, &product.Description, &product.Price, &product.Amount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &products.ErrNotFound{ProductId: id}
		}

		return nil, &errors.ErrInternal{Err: err}
	}

	return &product, nil
}

func (r *repository) FindByName(name string, ctx context.Context) (*products.Product, error) {
	var product products.Product

	err := r.pool.QueryRow(ctx, "SELECT product_id, name, description, price, amount FROM products WHERE name = $1", name).Scan(
		&product.Id, &product.Name, &product.Description, &product.Price, &product.Amount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &products.ErrNotFound{Name: name}
		}

		return nil, &errors.ErrInternal{Err: err}
	}

	return &product, nil
}

func (r *repository) List(ctx context.Context) ([]*products.Product, error) {
	var prods []*products.Product

	rows, err := r.pool.Query(ctx, "SELECT product_id, name, description, price, amount FROM products")
	if err != nil && err != pgx.ErrNoRows {
		return nil, &errors.ErrInternal{Err: err}
	}

	for rows.Next() {
		var p products.Product

		if err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Amount); err != nil {
			return nil, &errors.ErrInternal{Err: err}
		}

		prods = append(prods, &p)
	}

	return prods, nil
}

func (r *repository) Store(product *products.Product, ctx context.Context) error {
	if _, err := r.pool.Exec(
		ctx,
		"INSERT INTO products(product_id, name, description, price, amount) VALUES($1, $2, $3, $4, $5);",
		product.Id, product.Name, product.Description, product.Price, product.Amount,
	); err != nil {
		return err
	}

	return nil
}

func (*repository) Update(product *products.Product, ctx context.Context) error {
	panic("unimplemented")
}

func (r *repository) Delete(id products.ProductId, ctx context.Context) error {
	if _, err := r.pool.Exec(context.Background(), "DELETE FROM products WHERE product_id = $1;", id); err != nil {
		return err
	}

	return nil
}
