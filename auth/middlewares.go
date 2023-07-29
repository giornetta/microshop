package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/giornetta/microshop/errors"
	"github.com/giornetta/microshop/respond"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slices"
)

func AuthenticateMiddleware(verifier Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokens := strings.SplitAfter(r.Header.Get("Authorization"), "Bearer ")
			if len(tokens) < 2 {
				respond.Err(w, &errors.ErrUnauthorized{Err: fmt.Errorf("no token provided")})
				return
			}

			token, err := verifier.Verify(tokens[1])
			if err != nil {
				respond.Err(w, &errors.ErrUnauthorized{Err: err})
				return
			}

			ctx := context.WithValue(r.Context(), ContextKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthorizeSubjectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := FromContext(r.Context())
		if err != nil {
			respond.Err(w, &errors.ErrUnauthorized{Err: err})
			return
		}

		id := chi.URLParam(r, "id")
		if token.Subject() != id {
			respond.Err(w, &errors.ErrUnauthorized{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthorizeRolesMiddleware(allowedRoles []Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := FromContext(r.Context())
			if err != nil {
				respond.Err(w, &errors.ErrUnauthorized{Err: err})
				return
			}

			for _, role := range token.Roles() {
				if slices.Contains(allowedRoles, role) {
					next.ServeHTTP(w, r)
					return
				}
			}

			respond.Err(w, &errors.ErrUnauthorized{})
		})
	}
}
