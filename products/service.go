package products

import (
	"context"
	"log"

	"github.com/giornetta/microshop/events"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type service struct {
	Repository ProductRepository
	Producer   events.Publisher
	Validator  *validator.Validate
}

func NewService(repository ProductRepository, producer events.Publisher) Service {
	return &service{
		Repository: repository,
		Producer:   producer,
		Validator:  validator.New(),
	}
}

func (s *service) Create(req CreateProductRequest, ctx context.Context) (*Product, error) {
	if err := s.Validator.Struct(req); err != nil {
		return nil, err
	}

	product := &Product{
		Id:          ProductId(uuid.New().String()),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Amount:      req.InitialAmount,
	}

	if err := s.Producer.Publish(events.ProductCreated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx); err != nil {
		log.Println(err)
	}

	return product, nil
}

func (s *service) GetById(productId ProductId, ctx context.Context) (*Product, error) {
	return s.Repository.FindById(productId, ctx)
}

func (s *service) List(ctx context.Context) ([]*Product, error) {
	return s.Repository.List(ctx)
}

func (s *service) Update(req UpdateProductRequest, ctx context.Context) error {
	if err := s.Validator.Struct(req); err != nil {
		return err
	}

	product, err := s.Repository.FindById(req.Id, ctx)
	if err != nil {
		return nil
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
		return err
	}

	return nil
}

func (s *service) Restock(req RestockProductRequest, ctx context.Context) error {
	if err := s.Validator.Struct(req); err != nil {
		return err
	}

	product, err := s.Repository.FindById(req.Id, ctx)
	if err != nil {
		return nil
	}

	product.UpdateStock(req.Amount)

	if err := s.Producer.Publish(events.ProductUpdated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx); err != nil {
		return nil
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
		return err
	}

	return nil
}
