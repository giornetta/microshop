package identity

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/giornetta/microshop/auth"
	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/respond"
)

type handler struct {
	Service Service
}

func NewRouter(service Service) http.Handler {
	h := &handler{
		Service: service,
	}

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	router.Route("/api/v1/identity", func(r chi.Router) {
		r.Post("/signup", h.handleSignUp)
		r.Post("/signin", h.handleSignIn)
		r.Put("/roles/{id}", h.handleUpdateRoles)
	})

	return router
}

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var req signUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	p, err := h.Service.SignUp(&SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
	}, r.Context())
	if err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, p)
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var req signInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	p, err := h.Service.SignIn(&SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	}, r.Context())
	if err != nil {
		fmt.Println(err)
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, p)
}

type updateRolesRequest struct {
	Roles []auth.Role `json:"roles"`
}

func (h *handler) handleUpdateRoles(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Err(w, &errors.ErrBadRequest{})
		return
	}

	if err := h.Service.UpdateRoles(&UpdateRolesRequest{
		Id:    IdentityId(id),
		Roles: req.Roles,
	}, r.Context()); err != nil {
		respond.Err(w, err)
		return
	}

	respond.JSON(w, http.StatusNoContent, nil)
}
