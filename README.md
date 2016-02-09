# gobstopper

An idiomatic framework for creating web stacks in Golang.

[![Build Status](https://travis-ci.org/seedboxtech/gobstopper.svg?branch=master)](https://travis-ci.org/seedboxtech/gobstopper) [![](https://godoc.org/github.com//seedboxtech/gobstopper?status.svg)](http://godoc.org/github.com//seedboxtech/gobstopper)

**version** 3.0

## Goals

- Well chosen stack of core components
- Exposed core components for ultimate flexibility
- Convenience functions for adding routes, middleware and route groups for those looking to get going fast
- DB connection magically available in all handlers and middleware

## Usage

### Basic

```go
conf := gobstopper.Config{
  Port: 8000,
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

server.Start()
```

### Route Groups

```go

server.RouteGroup("/admin", func(group *gobstopper.Group) {
  group.Middleware(middleware.TerribleAuthMiddleware)
  group.Route(gobstopper.Route{
    Method: "GET",
    Path:   "",
    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintln(w, "admin")
    }),
  })
})


```

### Config

For full details of the server [configuration](https://godoc.org/github.com/seedboxtech/gobstopper#Config).



### Database Connection

Gobstopper will make your sqlx DB connection pool available to you in all your middleware and handlers.

```go
// Route handler...
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := database.GetConnection(r)
  // You can use it from here on...
}
```



## Under The Hood

### [mux](http://www.gorillatoolkit.org/pkg/mux) - Router

> The name mux stands for "HTTP request multiplexer".
> Like the standard http.ServeMux, mux.Router matches incoming
> requests against a list of registered routes and calls a handler
> for the route that matches the URL or other conditions.

### [negroni](https://github.com/codegangsta/negroni) - Middleware Manager

> Negroni is an idiomatic approach to web middleware in Go.
> It is tiny, non-intrusive, and encourages use of net/http Handlers.

### [sqlx](https://github.com/jmoiron/sqlx) - Database (optional)

> sqlx is a library which provides a set of extensions on go's standard database/sql library.
> The sqlx versions of sql.DB, sql.TX, sql.Stmt, et al.
> all leave the underlying interfaces untouched, so that
> their interfaces are a superset on the standard ones.


## Contributing

See an issue? Submit a [ticket](https://github.com/seedboxtech/gobstopper/issues).

This project aims to have unit tests with 100% code coverage. All PRs should have full tests included with them. To run the tests...

```bash
go test -v ./...
```
