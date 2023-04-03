package products

import (
	"fmt"
)

type ErrInternal struct {
	Err error
}

func (e *ErrInternal) Error() string {
	return "Internal Server Error"
}

func (e *ErrInternal) Cause() error {
	return e.Err
}

type ErrNotFound struct {
	ProductId ProductId
	Name      string
}

func (err *ErrNotFound) Error() string {
	if err.ProductId != "" {
		return fmt.Sprintf("product with id=%s was not found", err.ProductId.String())
	}

	if err.Name != "" {
		return fmt.Sprintf("product with name=%s was not found", err.Name)
	}

	return "product was not found"
}
