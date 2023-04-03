package products

import (
	"context"
	"fmt"

	"github.com/giornetta/microshop/events"
	"golang.org/x/exp/slog"
)

type productHandler struct {
	repository ProductRepository
}

func NewProductHandler(repository ProductRepository) events.Handler {
	return &productHandler{
		repository: repository,
	}
}

func (h *productHandler) Handle(evt events.Event, ctx context.Context) error {
	var err error

	switch evt.Type() {
	case events.ProductCreatedType:
		err = h.handleCreated(evt.(events.ProductCreated), ctx)
	case events.ProductUpdatedType:
		err = h.handleUpdated(evt.(events.ProductUpdated), ctx)
	case events.ProductDeletedType:
		err = h.handleDeleted(evt.(events.ProductDeleted), ctx)
	default:
		err = fmt.Errorf("unknown event type: %v", evt.Type())
	}

	return err
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

type loggingHandler struct {
	handler events.Handler
	logger  *slog.Logger
}

func NewLoggingEventHandler(logger *slog.Logger, handler events.Handler) events.Handler {
	return &loggingHandler{
		handler: handler,
		logger:  logger,
	}
}

func (h *loggingHandler) Handle(evt events.Event, ctx context.Context) error {
	if err := h.handler.Handle(evt, ctx); err != nil {
		if e, ok := err.(*ErrInternal); ok {
			h.logger.Error("could not handle event",
				slog.String("type", evt.Type().String()),
				slog.String("err", e.Cause().Error()))
		}

		return err
	}

	return nil
}
