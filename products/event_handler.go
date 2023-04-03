package products

import (
	"context"
	"fmt"

	"github.com/giornetta/microshop/events"
)

type productHandler struct {
	repository ProductRepository
}

func NewProductHandler(repository ProductRepository) events.Handler {
	return &productHandler{
		repository: repository,
	}
}

func (h *productHandler) Handle(e events.Event, ctx context.Context) error {
	switch e.Type() {
	case events.ProductCreatedType:
		return h.handleCreated(e.(events.ProductCreated), ctx)
	case events.ProductUpdatedType:
		return h.handleUpdated(e.(events.ProductUpdated), ctx)
	case events.ProductDeletedType:
		return h.handleDeleted(e.(events.ProductDeleted), ctx)
	default:
		return fmt.Errorf("unknown event type: %v", e.Type())
	}
}

func (h *productHandler) handleCreated(evt events.ProductCreated, ctx context.Context) error {
	p := &Product{
		Id:          ProductId(evt.ProductId),
		Name:        evt.Name,
		Description: evt.Description,
		Price:       evt.Price,
		Amount:      evt.Amount,
	}

	if err := h.repository.Store(p, ctx); err != nil {
		return err
	}

	return nil
}

func (h *productHandler) handleUpdated(evt events.ProductUpdated, ctx context.Context) error {
	p := &Product{
		Id:          ProductId(evt.ProductId),
		Name:        evt.Name,
		Description: evt.Description,
		Price:       evt.Price,
		Amount:      evt.Amount,
	}

	if err := h.repository.Update(p, ctx); err != nil {
		return err
	}

	return nil
}

func (h *productHandler) handleDeleted(evt events.ProductDeleted, ctx context.Context) error {
	if err := h.repository.Delete(ProductId(evt.ProductId), ctx); err != nil {
		return err
	}

	return nil
}
