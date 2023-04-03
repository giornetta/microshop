package router

import (
	"encoding/json"
	"net/http"

	"github.com/giornetta/microshop/respond"

	"github.com/giornetta/microshop/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type handler struct {
	Service products.Service
}

func New(service products.Service) http.Handler {
	h := &handler{
		Service: service,
	}

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	router.Route("/api/v1/products", func(r chi.Router) {
		r.Post("/", h.handleCreateProduct)
		r.Get("/", h.handleListProducts)
		r.Get("/{id}", h.handleGetProduct)
		r.Put("/{id}", h.handleUpdateProduct)
		r.Delete("/{id}", h.handleDeleteProduct)
	})

	return router
}

type createProductRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float32 `json:"price"`
	InitialAmount int     `json:"initialAmount"`
}

func (h *handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.Service.Create(products.CreateProductRequest{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		InitialAmount: req.InitialAmount,
	}, r.Context())
	if err != nil {
		respond.Err(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusCreated, p)
}

func (h *handler) handleListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.List(r.Context())
	if err != nil {
		respond.Err(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, products)
}

func (h *handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "id")

	product, err := h.Service.GetById(products.ProductId(productId), r.Context())
	if err != nil {
		respond.Err(w, http.StatusNotFound, err)
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
		respond.Err(w, http.StatusBadRequest, err)
		return
	}

	if err := h.Service.Update(products.UpdateProductRequest{
		Id:          products.ProductId(id),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}, r.Context()); err != nil {
		respond.Err(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}

func (h *handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "id")

	if err := h.Service.Delete(products.ProductId(productId), r.Context()); err != nil {
		respond.Err(w, http.StatusNotFound, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}
