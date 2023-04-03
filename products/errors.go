package products

import (
	"fmt"
	"net/http"
)

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

func (err *ErrNotFound) StatusCode() int {
	return http.StatusNotFound
}
