package products

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ProductId string

func (id ProductId) String() string {
	return string(id)
}

type Product struct {
	Id          ProductId `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Amount      int       `json:"amount"`
}

// TODO This should be changed
func (p *Product) UpdateStock(amountDelta int) {
	p.Amount += amountDelta
	if p.Amount < 0 {
		p.Amount = 0
	}
}

type ProductRepository interface {
	Store(product *Product, ctx context.Context) error
	FindById(id ProductId, ctx context.Context) (*Product, error)
	FindByName(name string, ctx context.Context) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
	Update(product *Product, ctx context.Context) error
	Delete(id ProductId, ctx context.Context) error
}

type Service interface {
	Create(req *CreateProductRequest, ctx context.Context) (*Product, error)
	GetById(productId ProductId, ctx context.Context) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
	Update(req *UpdateProductRequest, ctx context.Context) error
	Restock(req *RestockProductRequest, ctx context.Context) error
	Delete(productId ProductId, ctx context.Context) error
}

type CreateProductRequest struct {
	Name        string
	Description string
	Price       float32
	Amount      int
}

func (r *CreateProductRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)

	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required,
			validation.Length(4, 32),
			is.ASCII,
		),
		validation.Field(&r.Description,
			validation.Required,
			validation.Length(10, 256),
			is.ASCII,
		),
		validation.Field(&r.Price,
			validation.Required,
			validation.Min(0).Exclusive(),
		),
		validation.Field(&r.Amount,
			validation.Min(0),
		),
	)
}

type UpdateProductRequest struct {
	Id          ProductId
	Name        string
	Description string
	Price       float32
}

func (r *UpdateProductRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)

	return validation.ValidateStruct(r,
		validation.Field(&r.Id,
			validation.Required,
		),
		validation.Field(&r.Name,
			validation.Required.When(r.Description == "" && r.Price == 0),
			validation.Length(4, 32),
			is.ASCII,
		),
		validation.Field(&r.Description,
			validation.Required.When(r.Name == "" && r.Price == 0),
			validation.Length(10, 256),
			is.ASCII,
		),
		validation.Field(&r.Price,
			validation.Required.When(r.Description == "" && r.Name == ""),
			validation.Min(0.).Exclusive(),
		),
	)
}

type RestockProductRequest struct {
	Id     ProductId
	Amount int
}

func (r *RestockProductRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Id,
			validation.Required,
		),
		validation.Field(&r.Amount,
			validation.Required,
			validation.Min(0).Exclusive(),
		),
	)
}
