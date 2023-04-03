package products

import (
	"context"
)

type ProductId string

func (id ProductId) String() string {
	return string(id)
}

type Product struct {
	Id          ProductId `db:"product_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float32   `db:"price"`
	Amount      int       `db:"amount"`
}

// TODO This should be changed
func (p *Product) UpdateStock(amountDelta int) {
	p.Amount += amountDelta
	if p.Amount < 0 {
		p.Amount = 0
	}
}

type Service interface {
	Create(req CreateProductRequest, ctx context.Context) (*Product, error)
	GetById(productId ProductId, ctx context.Context) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
	Update(req UpdateProductRequest, ctx context.Context) error
	Restock(req RestockProductRequest, ctx context.Context) error
	Delete(productId ProductId, ctx context.Context) error
}

type CreateProductRequest struct {
	Name          string  `validate:"required,alphanum"`
	Description   string  `validate:"required,min=10"`
	Price         float32 `validate:"required,gt=0"`
	InitialAmount int     `validate:"gte=0"`
}

type UpdateProductRequest struct {
	Id          ProductId `validate:"required"`
	Name        string    `validate:"alphanum"`
	Description string    `validate:"min=10"`
	Price       float32   `validate:"gte=0"`
}

type RestockProductRequest struct {
	Id     ProductId `validate:"required"`
	Amount int       `validate:"required,gt=0"`
}

type ProductRepository interface {
	Store(product *Product, ctx context.Context) error
	FindById(id ProductId, ctx context.Context) (*Product, error)
	FindByName(name string, ctx context.Context) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
	Update(product *Product, ctx context.Context) error
	Delete(id ProductId, ctx context.Context) error
}
