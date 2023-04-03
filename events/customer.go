package events

const CustomerTopic Topic = "Customers"

func init() {

}

const (
	CustomerCreatedType                Type = "Customer.Created"
	CustomerShippingAddressUpdatedType Type = "Customer.ShippingAddressUpdated"
	CustomerDeletedType                Type = "Customer.Deleted"
)

type CustomerEvent struct {
	CustomerId string `json:"customer_id"`
}

func (e CustomerEvent) Key() Key { return Key(e.CustomerId) }

func (CustomerEvent) Topic() Topic { return CustomerTopic }

type CustomerCreated struct {
	CustomerEvent

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (CustomerCreated) Type() Type { return CustomerCreatedType }

type CustomerShippingAddressUpdated struct {
	CustomerEvent

	Country string `json:"country"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
	Street  string `json:"street"`
}

func (CustomerShippingAddressUpdated) Type() Type { return CustomerShippingAddressUpdatedType }

type CustomerDeleted struct {
	CustomerEvent
}

func (CustomerDeleted) Type() Type { return CustomerDeletedType }
