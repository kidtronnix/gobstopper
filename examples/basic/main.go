package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/seedboxtech/gobstopper"
	"github.com/seedboxtech/gobstopper/examples/middleware"
)

func main() {
	conf := gobstopper.Config{
		Port: 8000,
		// Negroni: negroni.New(),
		PathPrefix: "/v1",
	}

	server, err := gobstopper.NewServer(conf)
	if err != nil {
		log.Fatal("Unable to start server!")
	}

	server.Route(gobstopper.Route{
		Method: "GET",
		Path:   "/foo",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello")
		}),
	})

	server.RouteGroup("/admin", func(group *gobstopper.Group) {
		group.Middleware(middleware.TerribleAuthMiddleware{})
		group.Route(gobstopper.Route{
			Method: "GET",
			Path:   "",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "admin")
			}),
		})
	})

	server.Start()

}
