package lib

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Request struct {
	RawRequest       *http.Request
	RawBody          []byte
	Body             interface{}
	RequestPrincipal *Principal
	QueryParam       map[string]string
	Logger           *zerolog.Logger
}
