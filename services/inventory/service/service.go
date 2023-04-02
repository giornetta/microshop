package service

import (
	"context"
	"log"

	"github.com/giornetta/microshop/events"
	"github.com/giornetta/microshop/services/inventory"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
)

type service struct {
	Repository inventory.ProductRepository
	Producer   events.Publisher
	Validator  *validator.Validate
}

func New(repository inventory.ProductRepository, producer events.Publisher) inventory.ProductService {
	return &service{
		Repository: repository,
		Producer:   producer,
		Validator:  validator.New(),
	}
}

func (s *service) Create(req inventory.CreateProductRequest, ctx context.Context) (*inventory.Product, error) {
	if err := s.Validator.Struct(req); err != nil {
		return nil, err
	}

	product := &inventory.Product{
		Id:          inventory.ProductId(ksuid.New().String()),
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

func (s *service) GetById(productId inventory.ProductId, ctx context.Context) (*inventory.Product, error) {
	return s.Repository.FindById(productId)
}

func (s *service) List(ctx context.Context) ([]*inventory.Product, error) {
	return s.Repository.List()
}

func (s *service) Update(req inventory.UpdateProductRequest, ctx context.Context) error {
	if err := s.Validator.Struct(req); err != nil {
		return err
	}

	product, err := s.Repository.FindById(req.Id)
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

func (s *service) Restock(req inventory.RestockProductRequest, ctx context.Context) error {
	if err := s.Validator.Struct(req); err != nil {
		return err
	}

	product, err := s.Repository.FindById(req.Id)
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

func (s *service) Delete(productId inventory.ProductId, ctx context.Context) error {
	if _, err := s.Repository.FindById(productId); err != nil {
		return err
	}

	if err := s.Producer.Publish(events.ProductDeleted{
		ProductEvent: events.ProductEvent{ProductId: productId.String()},
	}, ctx); err != nil {
		return err
	}

	return nil
}
