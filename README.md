# gobstopper

A highly opinionated golang web stack.

[![Build Status](https://travis-ci.org/seedboxtech/gobstopper.svg?branch=master)](https://travis-ci.org/seedboxtech/gobstopper) [![](https://godoc.org/github.com//seedboxtech/gobstopper?status.svg)](http://godoc.org/github.com//seedboxtech/gobstopper)

**version** 2.1

## Stack

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

## Usage

Look at `main.go` to see a basic usage of a gobstopper stack. Here's an example of how to start it...

```bash
go run main.go -port 8000 -prefix "/v1" -connection "mysql|user:pass@tcp(host:3306)/db"
```

## Adding Middleware

To see an example of some custom middleware take a look at `middleware/example/example.go`. To see how to add it to the stack, see `main.go`.


### Core Packages

#### `github.com/seedboxtech/gobstopper/db` - [Go docs](https://godoc.org/github.com/seedboxtech/gobstopper/db)

#### `github.com/seedboxtech/gobstopper/service` - [Go docs](https://godoc.org/github.com/seedboxtech/gobstopper/service)

## Contributing

This project aims to have unit tests with 100% code coverage as well as a functional test of the main package. All PRs should have full tests included with them. To run the tests...

```bash
go test -v ./...
```
