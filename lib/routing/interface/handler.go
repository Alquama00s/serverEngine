package routingInterface

import (
	"net/http"

	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type Handler interface {
	GetHandler(func(*routingModel.Request) (*routingModel.Response, error)) func(http.ResponseWriter, *http.Request)
	GetRequestProcessors(path string) []RequestProcessor
	GetResponseProcessors(path string) []ResponseProcessor
	GetErrorHandler() *func(error, *routingModel.Request, *routingModel.Response) error
	SetErrorHandler(fun *func(error, *routingModel.Request, *routingModel.Response) error)
	AddRequestProcessor(regex string, priority int, processor RequestProcessor)
	AddResponseProcessor(regex string, priority int, processor ResponseProcessor)
	Finalize()
}
