package routingModel

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Request struct {
	RawRequest *http.Request
	RawBody    []byte
	Body       interface{}
	// RequestPrincipal *lib.Principal
	QueryParam map[string]string
	Logger     *zerolog.Logger
	metaData   map[string]interface{}
}

func (r *Request) GetMetaData(key string) interface{} {
	m, ok := r.metaData[key]
	if ok {
		return m
	}
	return nil
}

func (r *Request) SetMetaData(key string, value interface{}) {
	if r.metaData == nil {
		r.metaData = make(map[string]interface{})
	}
	r.metaData[key] = value
}
