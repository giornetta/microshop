package service

import (
	"fmt"

	"github.com/giornetta/microshop/events"
	"github.com/giornetta/microshop/services/inventory"
)

type productHandler struct {
	repository inventory.ProductRepository
}

func NewProductHandler(repository inventory.ProductRepository) events.Handler {
	return &productHandler{
		repository: repository,
	}
}

func (h *productHandler) Handle(e events.Event) error {
	switch e.Type() {
	case events.ProductCreatedType:
		return h.handleCreated(e.(events.ProductCreated))
	case events.ProductUpdatedType:
		return h.handleUpdated(e.(events.ProductUpdated))
	case events.ProductDeletedType:
		return h.handleDeleted(e.(events.ProductDeleted))
	default:
		return fmt.Errorf("unknown event type: %v", e.Type())
	}
}

func (h *productHandler) handleCreated(evt events.ProductCreated) error {
	p := &inventory.Product{
		Id:          inventory.ProductId(evt.ProductId),
		Name:        evt.Name,
		Description: evt.Description,
		Price:       evt.Price,
		Amount:      evt.Amount,
	}

	if err := h.repository.Store(p); err != nil {
		return err
	}

	return nil
}

func (h *productHandler) handleUpdated(evt events.ProductUpdated) error {
	p := &inventory.Product{
		Id:          inventory.ProductId(evt.ProductId),
		Name:        evt.Name,
		Description: evt.Description,
		Price:       evt.Price,
		Amount:      evt.Amount,
	}

	if err := h.repository.Update(p); err != nil {
		return err
	}

	return nil
}

func (h *productHandler) handleDeleted(evt events.ProductDeleted) error {
	if err := h.repository.Delete(inventory.ProductId(evt.ProductId)); err != nil {
		return err
	}

	return nil
}
