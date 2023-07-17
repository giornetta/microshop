package customers

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CustomerId string

func (id CustomerId) String() string {
	return string(id)
}

type Customer struct {
	ShippingAddress *ShippingAddress `json:"shipping_address"`
	Id              CustomerId       `json:"customer_id"`
	FirstName       string           `json:"first_name"`
	LastName        string           `json:"last_name"`
	Email           string           `json:"email"`
}

type ShippingAddress struct {
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	ZipCode string `json:"zip_code,omitempty"`
	Street  string `json:"street,omitempty"`
}

type CustomerQuerier interface {
	FindById(id CustomerId, ctx context.Context) (*Customer, error)
	FindByEmail(email string, ctx context.Context) (*Customer, error)
}

type CustomerStorer interface {
	Store(customer *Customer, ctx context.Context) error
	UpdateShippingAddress(id CustomerId, addr *ShippingAddress, ctx context.Context) error
	Delete(id CustomerId, ctx context.Context) error
}

type CustomerRepository interface {
	CustomerQuerier
	CustomerStorer
}

type Service interface {
	Create(req *CreateCustomerRequest, ctx context.Context) (*Customer, error)
	GetById(customerId CustomerId, ctx context.Context) (*Customer, error)
	UpdateShippingAddress(req *UpdateShippingAddressRequest, ctx context.Context) error
	Delete(customerId CustomerId, ctx context.Context) error
}

type CreateCustomerRequest struct {
	FirstName string
	LastName  string
	Email     string
}

func (r *CreateCustomerRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.FirstName,
			validation.Required,
			validation.Length(2, 10),
			is.Alpha,
		),
		validation.Field(&r.LastName,
			validation.Required,
			validation.Length(2, 10),
			is.Alpha,
		),
		validation.Field(&r.Email,
			validation.Required,
			is.EmailFormat,
		),
	)
}

type UpdateShippingAddressRequest struct {
	Id CustomerId

	Country string
	City    string
	ZipCode string
	Street  string
}

func (r *UpdateShippingAddressRequest) Validate() error {
	r.Country = strings.TrimSpace(r.Country)
	r.City = strings.TrimSpace(r.City)
	r.Street = strings.TrimSpace(r.Street)

	return validation.ValidateStruct(r,
		validation.Field(&r.Country,
			validation.Required,
			validation.Length(4, 12),
			is.ASCII,
		),
		validation.Field(&r.City,
			validation.Required,
			validation.Length(4, 12),
			is.Alpha,
		),
		validation.Field(&r.ZipCode,
			validation.Required,
			is.Int,
		),
		validation.Field(&r.Street,
			validation.Required,
			validation.Length(6, 32),
			is.ASCII,
		),
	)
}
