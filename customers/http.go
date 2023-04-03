package customers

import (
	"encoding/json"
	"net/http"

	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/respond"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type handler struct {
	service Service
}

func NewRouter(service Service) http.Handler {
	h := &handler{
		service: service,
	}

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	router.Route("/api/v1/customers", func(r chi.Router) {
		r.Post("/", h.handleCreateCustomer)
		r.Get("/{id}", h.handleGetCustomer)
		r.Put("/{id}/shipping", h.handleUpdateShippingAddress)
		r.Delete("/{id}", h.handleDeleteCustomer)
	})

	return router
}

type createCustomerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (h *handler) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req createCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	c, err := h.service.Create(&CreateCustomerRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}, r.Context())
	if err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, c)
}

func (h *handler) handleGetCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "id")

	c, err := h.service.GetById(CustomerId(customerId), r.Context())
	if err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, c)
}

type updateShippingAddressRequest struct {
	Country string `json:"country"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
	Street  string `json:"street"`
}

func (h *handler) handleUpdateShippingAddress(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "id")

	var req updateShippingAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	if err := h.service.UpdateShippingAddress(&UpdateShippingAddressRequest{
		Id:      CustomerId(customerId),
		Country: req.Country,
		City:    req.City,
		ZipCode: req.ZipCode,
		Street:  req.Street,
	}, r.Context()); err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}

func (h *handler) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "id")

	if err := h.service.Delete(CustomerId(customerId), r.Context()); err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}
