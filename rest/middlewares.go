// Defines middlewares used by the rest API handlers

package rest

import (
	"context"
	"fmt"
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

// required returns a middleware that requires some fields to be present
// in request body or query
func required(next http.HandlerFunc, fields ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, f := range fields {
			if r.FormValue(f) == "" {
				o := Res{
					Msg:    fmt.Sprintf("Field %s is required", f),
					Status: http.StatusBadRequest,
				}
				o.Send(w)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// auth is a middleware that ensures valid user JWT is present in Authorization
// Bearer token and verfies token signing using given secret and `ist`
func ensureAuthenticatedUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(a) != 2 || a[0] != "Bearer" {
			errorRes(w, "Invalid Authorization header. Authorization header needs Bearer token.", http.StatusUnauthorized)
			return
		}
		claims, err := account.ParseUserToken(a[1])
		if err != nil || claims.IssuerType != "user" {
			errorRes(w, "Invalid authorization token.", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxClaims, claims)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
