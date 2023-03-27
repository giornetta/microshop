package inventory

import (
	"context"
	"log"

	"github.com/giornetta/microshop/events"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
)

type service struct {
	Repository ProductRepository
	Producer   events.Publisher
	Validator  *validator.Validate
}

func (s *service) Create(req CreateProductRequest, ctx context.Context) (*Product, error) {
	if err := s.Validator.Struct(req); err != nil {
		return nil, err
	}

	product := &Product{
		Id:          ProductId(ksuid.New().String()),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Amount:      req.InitialAmount,
	}

	if err := s.Repository.Store(product); err != nil {
		return nil, err
	}

	// TODO ERRCHECK and retries.
	err := s.Producer.Publish(events.ProductCreated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx)
	if err != nil {
		log.Println(err)
	}

	return product, nil
}

func (s *service) GetById(productId ProductId, ctx context.Context) (*Product, error) {
	return s.Repository.FindById(productId)
}

func (s *service) List(ctx context.Context) ([]*Product, error) {
	return s.Repository.List()
}

func (s *service) Update(req UpdateProductRequest, ctx context.Context) error {
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

	if err := s.Repository.Update(product); err != nil {
		return err
	}

	// TODO Errcheck and retries
	_ = s.Producer.Publish(events.ProductUpdated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx)

	return nil
}

func (s *service) Restock(req RestockProductRequest, ctx context.Context) error {
	if err := s.Validator.Struct(req); err != nil {
		return err
	}

	product, err := s.Repository.FindById(req.Id)
	if err != nil {
		return nil
	}

	product.Amount += req.Amount

	if err := s.Repository.Update(product); err != nil {
		return err
	}

	// TODO Errcheck and retries
	_ = s.Producer.Publish(events.ProductUpdated{
		ProductEvent: events.ProductEvent{ProductId: product.Id.String()},
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Amount:       product.Amount,
	}, ctx)

	return nil
}

func (s *service) Delete(productId ProductId, ctx context.Context) error {
	if err := s.Repository.Delete(productId); err != nil {
		return err
	}

	// TODO Errcheck and retries
	_ = s.Producer.Publish(events.ProductDeleted{
		ProductEvent: events.ProductEvent{ProductId: productId.String()},
	}, ctx)

	return nil
}

func NewService(repository ProductRepository, producer events.Publisher) ProductService {
	return &service{
		Repository: repository,
		Producer:   producer,
		Validator:  validator.New(),
	}
}
