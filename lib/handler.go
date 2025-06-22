package lib

import "net/http"

type Handler interface {
	GetHandler(func(*Request) (*Response, error)) func(http.ResponseWriter, *http.Request)
	GetRequestProcessors(path string) []RequestProcessor
	GetResponseProcessors(path string) []ResponseProcessor
	GetErrorHandler() *func(error, *Request, *Response) error
	SetErrorHandler(fun *func(error, *Request, *Response) error)
	AddRequestProcessor(regex string, priority int, processor RequestProcessor)
	AddResponseProcessor(regex string, priority int, processor ResponseProcessor)
	Finalize()
}
