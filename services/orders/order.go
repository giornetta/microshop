package orders

type OrderId string

type Status byte

const (
	StatusCreated Status = iota
	StatusAccepted
	StatusShipped
	StatusDelivered
)

type Order struct {
	OrderId OrderId
	Status  Status
}

type LineItem struct {
}
