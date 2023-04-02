package inventory

import (
	"context"
)

type ProductId string

func (id ProductId) String() string {
	return string(id)
}

type Product struct {
	Id          ProductId
	Name        string
	Description string
	Price       float32
	Amount      int
}

// TODO This should be changed
func (p *Product) UpdateStock(amountDelta int) {
	p.Amount += amountDelta
	if p.Amount < 0 {
		p.Amount = 0
	}
}

type ProductService interface {
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
	Store(product *Product) error
	FindById(id ProductId) (*Product, error)
	FindByName(name string) (*Product, error)
	List() ([]*Product, error)
	Update(product *Product) error
	Delete(id ProductId) error
}
