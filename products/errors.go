package products

import "fmt"

type ErrNotFound struct {
	ProductId ProductId
}

func (err *ErrNotFound) Error() string {
	return fmt.Sprintf("product with id=%s was not found", err.ProductId.String())
}

type ErrIDAlreadyExists struct {
	ProductId ProductId
}

func (err *ErrIDAlreadyExists) Error() string {
	return fmt.Sprintf("product with id=%s already exists", err.ProductId.String())
}
