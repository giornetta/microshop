package products

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/giornetta/microshop/auth"
	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/respond"
)

type handler struct {
	service  Service
	verifier auth.Verifier
}

func NewRouter(service Service, verifier auth.Verifier) http.Handler {
	h := &handler{
		service:  service,
		verifier: verifier,
	}

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	router.Route("/api/v1/products", func(r chi.Router) {
		r.Get("/", h.handleListProducts)
		r.Get("/{id}", h.handleGetProduct)

		protectedRoutes := r.With(
			auth.AuthenticateMiddleware(h.verifier),
			auth.AuthorizeRolesMiddleware([]auth.Role{auth.AdminRole}),
		)
		protectedRoutes.Post("/", h.handleCreateProduct)
		protectedRoutes.Put("/{id}", h.handleUpdateProduct)
		protectedRoutes.Put("/restock/{id}", h.handleRestockProduct)
		protectedRoutes.Delete("/{id}", h.handleDeleteProduct)
	})

	return router
}

type createProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Amount      int     `json:"amount"`
}

func (h *handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	p, err := h.service.Create(&CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Amount:      req.Amount,
	}, r.Context())
	if err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, p)
}

func (h *handler) handleListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.List(r.Context())
	if err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, products)
}

func (h *handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "id")

	product, err := h.service.GetById(ProductId(productId), r.Context())
	if err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, product)
}

type updateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

func (h *handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	if err := h.service.Update(&UpdateProductRequest{
		Id:          ProductId(id),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}, r.Context()); err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}

type restockProductRequest struct {
	Amount uint `json:"amount"`
}

func (h *handler) handleRestockProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req restockProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	if err := h.service.Restock(&RestockProductRequest{
		Id:     ProductId(id),
		Amount: int(req.Amount),
	}, r.Context()); err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}

func (h *handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "id")

	if err := h.service.Delete(ProductId(productId), r.Context()); err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}
