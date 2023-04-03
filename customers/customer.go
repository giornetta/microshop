package customers

import "context"

type CustomerId string

func (id CustomerId) String() string {
	return string(id)
}

type Customer struct {
	Id CustomerId `json:"customer_id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`

	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

type ShippingAddress struct {
	Country string `json:"country"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
	Street  string `json:"street"`
}

type Service interface {
	Create(req *CreateCustomerRequest, ctx context.Context) error
	GetById(customerId CustomerId, ctx context.Context) (*Customer, error)
	UpdateShippingAddress(req *UpdateShippingAddressRequest, ctx context.Context) error
	Delete(customerId CustomerId, ctx context.Context) error
}

type CreateCustomerRequest struct {
	FirstName string
	LastName  string
	Email     string
}

type UpdateShippingAddressRequest struct {
	Country string
	City    string
	ZipCode string
	Street  string
}

type CustomerRepository interface {
	Store(customer *Customer, ctx context.Context) error
	FindById(id CustomerId, ctx context.Context) (*Customer, error)
	FindByEmail(email string, ctx context.Context) (*Customer, error)
	Update(customer *Customer, ctx context.Context) error
	Delete(id CustomerId, ctx context.Context) error
}
