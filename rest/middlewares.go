// Defines middlewares used by the rest API handlers

package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/applait/xplex-rig/account"
)

// ctxKey is a private key type to share context values in requests
type ctxKey int

// Context key definitions
const (
	ctxClaims ctxKey = iota + 1
)

// ensureContentType determines that POST, PUT and PATCH requests use `Content-Type: application/json`
func ensureContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") {
			h.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("Content-Type") == "application/json" {
			h.ServeHTTP(w, r)
			return
		}
		ErrUnsupportedMediaType.Send(w)
	})
}

// auth is a middleware that ensures valid user JWT is present in Authorization
// Bearer token and verfies token signing using given secret and `ist`
func ensureAuthenticatedUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(a) != 2 || a[0] != "Bearer" {
			e := ErrInvalidCredentials
			e.Message = "Invalid Authorization header. Authorization header needs Bearer token"
			e.Send(w)
			return
		}
		claims, err := account.ParseUserToken(a[1])
		if err != nil || claims.IssuerType != "user" {
			ErrInvalidCredentials.Send(w)
			return
		}
		ctx := context.WithValue(r.Context(), ctxClaims, claims)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Triggered when no matching routes are found
var notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ErrNotFound.Send(w)
})

// Triggered when method is not allowed on handler
var methodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ErrMethodNotAllowed.Send(w)
})
