package routingModel

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Response struct {
	StatusCode int
	Headers    *http.Header
	Request    *Request
	Body       interface{}
	Logger     *zerolog.Logger
}

func NewRestResponse() *Response {
	hdr := &http.Header{}
	hdr.Add("Content-Type", "application/json")
	return &Response{
		Headers: hdr,
	}
}

func (r *Response) SetBody(body interface{}) *Response {
	r.Body = body
	return r
}

func (r *Response) AddToBody(key string, value interface{}) *Response {
	if r.Body == nil {
		r.Body = make(map[string]interface{})
	}
	r.Body.(map[string]interface{})[key] = value
	return r
}

func (r *Response) SetStatus(stats int) *Response {
	r.StatusCode = stats
	return r
}
