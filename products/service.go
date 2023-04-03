package products

import (
	"context"
	"log"

	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/events"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type service struct {
	Repository ProductRepository
	Producer   events.Publisher
}

func NewService(repository ProductRepository, producer events.Publisher) Service {
	return &service{
		Repository: repository,
		Producer:   producer,
	}
}

func (s *service) Create(req *CreateProductRequest, ctx context.Context) (*Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	_, err := s.Repository.FindByName(req.Name, ctx)
	if err == nil {
		return nil, &ErrAlreadyExists{Name: req.Name}
	}
	if err != nil {
		if _, ok := err.(*ErrNotFound); !ok {
			return nil, err
		}
	}

	product := &Product{
		Id:          ProductId(uuid.New().String()),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Amount:      req.Amount,
	}

	if err := s.Producer.Publish(events.ProductCreated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx); err != nil {
		return nil, &errors.ErrInternal{Err: err}
	}

	return product, nil
}

func (s *service) GetById(productId ProductId, ctx context.Context) (*Product, error) {
	p, err := s.Repository.FindById(productId, ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) List(ctx context.Context) ([]*Product, error) {
	prods, err := s.Repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return prods, nil
}

func (s *service) Update(req *UpdateProductRequest, ctx context.Context) error {
	if err := req.Validate(); err != nil {
		log.Println(err.Error())
		return &errors.ErrBadRequest{Err: err}
	}

	product, err := s.Repository.FindById(req.Id, ctx)
	if err != nil {
		return err
	}

	if req.Name != "" {
		product.Name = req.Name
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.Price != 0 {
		product.Price = req.Price
	}

	if err := s.Producer.Publish(events.ProductUpdated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx); err != nil {
		return &errors.ErrInternal{Err: err}
	}

	return nil
}

func (s *service) Restock(req *RestockProductRequest, ctx context.Context) error {
	if err := req.Validate(); err != nil {
		return err
	}

	product, err := s.Repository.FindById(req.Id, ctx)
	if err != nil {
		return err
	}

	product.UpdateStock(req.Amount)

	if err := s.Producer.Publish(events.ProductUpdated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx); err != nil {
		return &errors.ErrInternal{Err: err}
	}

	return nil
}

func (s *service) Delete(productId ProductId, ctx context.Context) error {
	if _, err := s.Repository.FindById(productId, ctx); err != nil {
		return err
	}

	if err := s.Producer.Publish(events.ProductDeleted{
		ProductEvent: events.ProductEvent{ProductId: productId.String()},
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

func (s *loggingService) Create(req *CreateProductRequest, ctx context.Context) (*Product, error) {
	p, err := s.service.Create(req, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not create product",
				slog.String("method", "Create"),
				slog.String("err", e.Cause().Error()),
			)
		}

		return nil, err
	}

	return p, nil
}

func (s *loggingService) GetById(productId ProductId, ctx context.Context) (*Product, error) {
	p, err := s.service.GetById(productId, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not find product by id",
				slog.String("method", "GetById"),
				slog.String("product_id", productId.String()),
				slog.String("err", e.Cause().Error()),
			)
		}

		return nil, err
	}

	return p, nil
}

func (s *loggingService) List(ctx context.Context) ([]*Product, error) {
	prods, err := s.service.List(ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not list products",
				slog.String("method", "List"),
				slog.String("err", e.Cause().Error()),
			)
		}

		return nil, err
	}

	return prods, nil
}

func (s *loggingService) Restock(req *RestockProductRequest, ctx context.Context) error {
	err := s.service.Restock(req, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not restock product",
				slog.String("method", "Restock"),
				slog.String("product_id", req.Id.String()),
				slog.String("err", e.Cause().Error()),
			)
		}

		return err
	}

	return nil
}

func (s *loggingService) Update(req *UpdateProductRequest, ctx context.Context) error {
	err := s.service.Update(req, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not update product",
				slog.String("method", "Update"),
				slog.String("product_id", req.Id.String()),
				slog.String("err", e.Cause().Error()),
			)
		}

		return err
	}

	return nil
}

func (s *loggingService) Delete(productId ProductId, ctx context.Context) error {
	err := s.service.Delete(productId, ctx)
	if err != nil {
		if e, ok := err.(*errors.ErrInternal); ok {
			s.logger.Error("could not delete product",
				slog.String("method", "Delete"),
				slog.String("product_id", productId.String()),
				slog.String("err", e.Cause().Error()),
			)
		}

		return err
	}

	return nil
}
