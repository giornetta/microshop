package customers

import (
	"context"

	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/events"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type service struct {
	querier   CustomerQuerier
	publisher events.Publisher
}

func NewService(querier CustomerQuerier, publisher events.Publisher) Service {
	return &service{
		querier:   querier,
		publisher: publisher,
	}
}

func (s *service) Create(req *CreateCustomerRequest, ctx context.Context) (*Customer, error) {
	if err := req.Validate(); err != nil {
		return nil, &errors.ErrBadRequest{Err: err}
	}

	_, err := s.querier.FindByEmail(req.Email, ctx)
	if err == nil {
		return nil, &ErrAlreadyExists{Email: req.Email}
	}
	if err != nil {
		if _, ok := err.(*ErrNotFound); !ok {
			return nil, err
		}
	}

	c := &Customer{
		Id:        CustomerId(uuid.NewString()),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if err := s.publisher.Publish(&events.CustomerCreated{
		CustomerEvent: events.CustomerEvent{CustomerId: c.Id.String()},
		FirstName:     c.FirstName,
		LastName:      c.LastName,
		Email:         c.Email,
	}, ctx); err != nil {
		return nil, &errors.ErrInternal{Err: err}
	}

	return c, nil
}

func (s *service) GetById(customerId CustomerId, ctx context.Context) (*Customer, error) {
	c, err := s.querier.FindById(customerId, ctx)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *service) UpdateShippingAddress(req *UpdateShippingAddressRequest, ctx context.Context) error {
	if err := req.Validate(); err != nil {
		return &errors.ErrBadRequest{Err: err}
	}

	_, err := s.querier.FindById(req.Id, ctx)
	if err != nil {
		return err
	}

	if err := s.publisher.Publish(events.CustomerShippingAddressUpdated{
		CustomerEvent: events.CustomerEvent{CustomerId: req.Id.String()},
		Country:       req.Country,
		City:          req.City,
		ZipCode:       req.ZipCode,
		Street:        req.Street,
	}, ctx); err != nil {
		return &errors.ErrInternal{Err: err}
	}

	return nil
}

func (s *service) Delete(customerId CustomerId, ctx context.Context) error {
	if _, err := s.querier.FindById(customerId, ctx); err != nil {
		return err
	}

	if err := s.publisher.Publish(events.CustomerDeleted{
		CustomerEvent: events.CustomerEvent{CustomerId: customerId.String()},
	}, ctx); err != nil {
		return &errors.ErrInternal{Err: err}
	}

	return nil
}

type loggingService struct {
	service Service
	logger  *slog.Logger
}

func NewLoggingService(logger *slog.Logger, service Service) Service {
	return &loggingService{
		service: service,
		logger:  logger,
	}
}

func (s *loggingService) Create(req *CreateCustomerRequest, ctx context.Context) (*Customer, error) {
	p, err := s.service.Create(req, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not create customer",
				slog.String("method", "Create"),
				slog.String("err", e.Cause().Error()),
			)
		}

		return nil, err
	}

	return p, nil
}

func (s *loggingService) GetById(customerId CustomerId, ctx context.Context) (*Customer, error) {
	p, err := s.service.GetById(customerId, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not find customer by id",
				slog.String("method", "GetById"),
				slog.String("customer_id", customerId.String()),
				slog.String("err", e.Cause().Error()),
			)
		}

		return nil, err
	}

	return p, nil
}

func (s *loggingService) UpdateShippingAddress(req *UpdateShippingAddressRequest, ctx context.Context) error {
	err := s.service.UpdateShippingAddress(req, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not update shipping address",
				slog.String("method", "UpdateShippingAddress"),
				slog.String("customer_id", req.Id.String()),
				slog.String("err", e.Cause().Error()),
			)
		}

		return err
	}

	return nil
}

func (s *loggingService) Delete(customerId CustomerId, ctx context.Context) error {
	err := s.service.Delete(customerId, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not delete customer",
				slog.String("method", "Delete"),
				slog.String("customer_id", customerId.String()),
				slog.String("err", e.Cause().Error()),
			)
		}

		return err
	}

	return nil
}
