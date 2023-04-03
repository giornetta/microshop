package log

import (
	"context"

	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/events"
	"golang.org/x/exp/slog"
)

type loggingHandler struct {
	handler events.Handler
	logger  *slog.Logger
}

func NewEventHandler(logger *slog.Logger, handler events.Handler) events.Handler {
	return &loggingHandler{
		handler: handler,
		logger:  logger,
	}
}

func (h *loggingHandler) Handle(evt events.Event, ctx context.Context) error {
	if err := h.handler.Handle(evt, ctx); err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			h.logger.Error("could not handle event",
				slog.String("type", evt.Type().String()),
				slog.String("err", e.Cause().Error()))
		}

		return err
	}

	return nil
}
