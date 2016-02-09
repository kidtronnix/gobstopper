package middleware

import (
	"fmt"
	"net/http"
)

// TerribleAuthMiddleware is an example of some middleware for gobstopper. It satisfies the negroni.Handler interface.
type TerribleAuthMiddleware struct{}

func (t TerribleAuthMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.FormValue("token")
	fmt.Println("Performing Authentication check...")
	if token == "123" {
		// Yield to the next request handler
		next(rw, r)
		// Can do stuff after if you wanna
		return
	}
	rw.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(rw, "Authentication failed")
}
