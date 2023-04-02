package inmem

import (
	"sync"

	"github.com/giornetta/microshop/services/inventory"
)

type repository struct {
	Products map[inventory.ProductId]*inventory.Product
	Mutex    sync.RWMutex
}

func NewRepository() inventory.ProductRepository {
	return &repository{
		Products: make(map[inventory.ProductId]*inventory.Product),
		Mutex:    sync.RWMutex{},
	}
}

func (r *repository) Store(product *inventory.Product) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Products[product.Id] = product

	return nil
}

func (r *repository) FindById(id inventory.ProductId) (*inventory.Product, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	p, ok := r.Products[id]
	if !ok {
		return nil, &inventory.ErrNotFound{ProductId: id}
	}

	return p, nil
}

func (r *repository) FindByName(name string) (*inventory.Product, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for _, p := range r.Products {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, &inventory.ErrNotFound{}
}

func (r *repository) List() ([]*inventory.Product, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	products := make([]*inventory.Product, 0, len(r.Products))
	for _, p := range r.Products {
		products = append(products, p)
	}

	return products, nil
}

func (r *repository) Update(product *inventory.Product) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	_, ok := r.Products[product.Id]
	if !ok {
		return &inventory.ErrNotFound{ProductId: product.Id}
	}

	r.Products[product.Id] = product

	return nil
}

func (r *repository) Delete(id inventory.ProductId) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if _, ok := r.Products[id]; !ok {
		return &inventory.ErrNotFound{ProductId: id}
	}

	delete(r.Products, id)

	return nil
}
