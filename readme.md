# Server Engine

A small, lightweight Go HTTP "server engine" providing a simple routing layer, request/response objects,
middleware (request & response processors), authentication adapters (Basic / JWT), and a registrar to
organize routes and initializers.

This project is organized as a small framework you can embed into your own `main` application to
register controllers, initializers and start an HTTP server with minimal boilerplate.

## Highlights

- Simple Router / Route DSL: create routers and attach routes with method + controller function
- Centralized Registrar to collect routers, controller sets and initializers
- Pluggable Handler (default: `SimpleHandler`) with request/response middleware support
- Built-in authenticators: Basic and JWT with a delegating authenticator helper
- Route-level authorization helper `SimpleRouteAuthorizer`
- Small, testable building blocks (Request/Response/Principal, middleware interfaces)

## Table of contents

- [Quick start](#quick-start)
- [Core concepts](#core-concepts)
- [Examples](#examples)
- [API pointers](#api-pointers)
- [Development & run](#development--run)
- [Next steps](#next-steps)

## Quick start

Requirements: Go 1.20+ (module-enabled project in this repo).

There is no `main` package included by default. Create a small `main.go` that bootstraps the
registrar and starts the server. Example:

```go
package main

import (
		"github.com/Alquama00s/serverEngine"
		"github.com/Alquama00s/serverEngine/lib"
		"net/http"
)

// simple controller set implementation
type MyControllers struct{}

func (m *MyControllers) Controllers() {
		// register routes using the global Registrar()
		r := serverEngine.Registrar()
		api := r.Router("/api")

		api.Path("/hello").Method("GET").Handeler(func(req *lib.Request) (*lib.Response, error) {
				return lib.NewRestResponse().SetBody(map[string]string{"hello": "world"}).SetStatus(http.StatusOK), nil
		})
}

func main() {
		// register your ControllerSet
		serverEngine.Registrar().RegisterControllerSet(&MyControllers{})

		// start the server (calls Initialize / Finalize internally)
		serverEngine.Sereve()
}
```

Then run:

```bash
go run main.go
```

Note: the package exposes `serverEngine.Sereve()` (yes, spelled that way in the repo) which initializes
registered initializers and finalizes routes, then starts an HTTP server on port `:8080`.

## Core concepts

- Router (`lib.Router`): a grouping with a `PathPrefix` and a list of `Route`s. Create one via
	`serverEngine.Registrar().Router("/prefix")`.

- Route (`lib.Route`): represents a single endpoint. Fluent methods include:
	- `Path(string)` to create a route within a router
	- `Method(string)` to set the HTTP method (e.g. "GET", "POST")
	- `Handeler(func(*lib.Request)(*lib.Response,error))` to set the controller function

- Registrar (`registrar.DefaultRegistrar`): global registry returned by `serverEngine.Registrar()`.
	Use it to register `ControllerSet`s and `Initializers`, and to create routers.

- Handler (`lib.Handler`, default implementation: `lib.SimpleHandler`): converts controller
	functions into `http.HandlerFunc`s, runs configured request/response processors, and serializes
	`lib.Response` bodies to JSON.

- Request / Response (`lib.Request`, `lib.Response`): typed wrappers around the raw `http.Request`
	and response payload. `Response` includes headers, status and body helpers like `NewRestResponse()`.

- Middleware
	- Request processors implement `lib.RequestProcessor` (example implementation: `lib.SimpleReqMiddleWare`)
	- Response processors implement `lib.ResponseProcessor` (example implementation: `lib.SimpleResMiddleWare`)
	Handlers select middleware by regex against the request path and run them in priority order.

- Authenticators
	- `lib.BasicAuthenticator` parses Basic auth headers
	- `lib.JWTAuthenticator` creates & validates JWTs (helpers: `CreateToken`, `ParsePrincipal`)
	- `lib.DelegatingAuthenticator` delegates to the right authenticator based on the header token type

- Authorization
	- `lib.SimpleRouteAuthorizer` exposes a helper that returns a `RequestProcessor` enforcing
		token type, roles and privileges for matching paths.

## Examples

1) Registering a route (fluent API)

```go
r := serverEngine.Registrar()
api := r.Router("/v1")

api.Path("/ping").Method("GET").Handeler(func(req *lib.Request) (*lib.Response, error) {
		return lib.NewRestResponse().SetBody(map[string]string{"pong": "ok"}).SetStatus(200), nil
})
```

2) Adding a request middleware (e.g., a simple logger or JSON body parser)

```go
sh := &lib.SimpleHandler{}
rp := &lib.SimpleReqMiddleWare{}
rp.SetRegex("^/v1/.*")
rp.SetPriority(10)
rp.Process(func(r *lib.Request) (*lib.Request, error, *lib.Response) {
		// do something with the request, set r.Body or r.RequestPrincipal etc.
		return r, nil, nil
})
sh.AddRequestProcessor("^/v1/.*", 10, rp)
```

3) Route-level authorization using `SimpleRouteAuthorizer`

```go
authz := lib.NewSimpleRouteAuth().Path("^/v1/secure/.*").Privileges("read:secure").TokenType("Bearer")
serverEngine.Registrar().Handler.AddRequestProcessor(authz.GetRequestProcessor().GetRegexString(), -2147483648, authz.GetRequestProcessor())
# Or register it via the registrar/controller setup flow so it gets applied by the default handler
```

4) Creating a JWT using the built-in authenticator

```go
jwtAuth := lib.NewJwtAuthenticator()
token, err := jwtAuth.CreateToken([]string{"read"}, []string{"user"}, 123, "alice")
if err != nil { /* handle error */ }
// send token back to client as Bearer token
```

## API pointers (important files)

- `server_engine.go` — Registrar bootstrap and `Sereve()` function that starts the server
- `registrar/default_registrar.go` — collects routers, controller sets and initializers
- `lib/router.go`, `lib/route.go` — router & route DSL
- `lib/handler.go`, `lib/simple_handler.go` — Handler interface and default implementation
- `lib/request.go`, `lib/response.go` — typed Request/Response structures
- `lib/*_mw.go` — middleware primitives (`SimpleReqMiddleWare`, `SimpleResMiddleWare`)
- `lib/authenticator.go`, `lib/JWTAuthenticator.go`, `lib/delegatingAuthenticator.go` — authenticators
- `lib/routeAuthorizer.go` — helper for route-based authorization

## Development & run

Install dependencies and tidy modules:

```bash
go mod tidy
```

Run a small local program (example `main.go` shown above):

```bash
go run main.go
```

The server listens on port 8080 by default.

## Notes & gotchas

- The project's registrar pattern registers routes and finalizes them into an `http.ServeMux`.
- Middlewares are selected by regex on the path and ordered by integer `Priority`.
- `JWTAuthenticator` generates keys in memory (for dev/demo). For production you should
	persist keys or integrate with a proper JWK provider.
- The exported function to start the server is named `Sereve()` in the repository. If you prefer
	a different startup approach, create your own `main` that uses the registrar and handler directly.

## Next steps / ideas

- Add a `main.go` example in `examples/` to demonstrate a complete app flow
- Add unit tests around the handler middleware ordering and authenticator parsing
- Add configuration for bind address, TLS, and key storage for JWT auth

## License

No license specified in the repository. Add a license file if you intend to publish this project.

---