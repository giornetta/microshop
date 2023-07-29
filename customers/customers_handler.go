package customers

import (
	"context"
	"fmt"

	"github.com/giornetta/microshop/events"
)

type productHandler struct {
	repository CustomerRepository
}

func NewCustomerHandler(repository CustomerRepository) events.Handler {
	return &productHandler{
		repository: repository,
	}
}

func (h *productHandler) Handle(evt events.Event, ctx context.Context) error {
	var err error

	switch evt.Type() {
	case events.CustomerCreatedType:
		err = h.handleCreated(evt.(events.CustomerCreated), ctx)
	case events.CustomerShippingAddressUpdatedType:
		err = h.handleShippingAddressUpdated(evt.(events.CustomerShippingAddressUpdated), ctx)
	case events.CustomerDeletedType:
		err = h.handleDeleted(evt.(events.CustomerDeleted), ctx)
	default:
		err = fmt.Errorf("unknown event type: %v", evt.Type())
	}

	return err
}

func (h *productHandler) handleCreated(evt events.CustomerCreated, ctx context.Context) error {
	p := &Customer{
		Id:        CustomerId(evt.CustomerId),
		FirstName: evt.FirstName,
		LastName:  evt.LastName,
	}

	if err := h.repository.Store(p, ctx); err != nil {
		return err
	}

	return nil
}

func (h *productHandler) handleShippingAddressUpdated(evt events.CustomerShippingAddressUpdated, ctx context.Context) error {
	addr := &ShippingAddress{
		Country: evt.Country,
		City:    evt.City,
		ZipCode: evt.ZipCode,
		Street:  evt.Street,
	}

	if err := h.repository.UpdateShippingAddress(CustomerId(evt.CustomerId), addr, ctx); err != nil {
		return err
	}

	return nil
}

func (h *productHandler) handleDeleted(evt events.CustomerDeleted, ctx context.Context) error {
	if err := h.repository.Delete(CustomerId(evt.CustomerId), ctx); err != nil {
		return err
	}

	return nil
}
