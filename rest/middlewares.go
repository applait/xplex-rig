// Defines middlewares used by the rest API handlers

package rest

import (
	"fmt"
	"net/http"
)

// required returns a middleware that requires some fields to be present
// in request body or query
func required(fields ...string) middleware {
	return func(h http.Handler) http.Handler {
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
			h.ServeHTTP(w, r)
		})
	}
}
