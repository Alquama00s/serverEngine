package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	routingI "github.com/Alquama00s/serverEngine/lib/routing/interface"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"

	"github.com/Alquama00s/serverEngine/lib/logging/loggerFactory"
)

var (
	logger = loggerFactory.GetLogger("SimpleHandler")
)

type SimpleHandler struct {
	requestMiddlewares  []routingI.RequestProcessor
	responseMiddlewares []routingI.ResponseProcessor
	errorHandler        *func(error, *routingModel.Request, *routingModel.Response) error
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	byt, err := json.Marshal(routingModel.NewError(err))
	if err != nil {
		fmt.Printf("error: %v", err)
		w.Write([]byte("error marshling"))
	} else {
		w.Write(byt)
	}
}

func (sh *SimpleHandler) handleResponse(w http.ResponseWriter, req *routingModel.Request, res *routingModel.Response) {
	byt, err := json.Marshal(res.Body)
	if err != nil {
		if sh.errorHandler != nil {
			err = (*sh.GetErrorHandler())(err, req, nil)
		}
		handleError(w, err)
		return
	}
	if res.StatusCode == 0 {
		res.StatusCode = http.StatusAccepted
	}
	for key, values := range *res.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.StatusCode)
	w.Write(byt)
}

func (sh *SimpleHandler) GetHandler(controller func(*routingModel.Request) (*routingModel.Response, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lgr := loggerFactory.GetLogger(r.URL.Path, r.Host)
		req := &routingModel.Request{
			RawRequest: r,
			Logger:     lgr,
		}
		var processingError error
		var res *routingModel.Response

		var tempRes *routingModel.Response
		var tempReq *routingModel.Request

		for _, rp := range sh.GetRequestProcessors(r.URL.Path) {
			tempReq, processingError, res = rp.GetProcessor()(req)
			if tempReq != nil {
				req = tempReq
			}
			if processingError != nil {
				if sh.errorHandler != nil {
					processingError = (*sh.GetErrorHandler())(processingError, req, nil)
				}
				handleError(w, processingError)
				return
			}
			if res != nil {
				res.Logger = lgr
				sh.handleResponse(w, req, res)
				return
			}
		}

		res, err := controller(req)

		if err != nil {
			if sh.errorHandler != nil {
				err = (*sh.GetErrorHandler())(err, req, res)
			}
			handleError(w, err)
			return
		}
		res.Logger = req.Logger
		res.Request = req
		for _, rp := range sh.GetResponseProcessors(r.URL.Path) {
			tempRes, processingError = rp.GetProcessor()(res)
			if tempRes != nil {
				res = tempRes
			}
			if processingError != nil {
				if sh.errorHandler != nil {
					processingError = (*sh.GetErrorHandler())(processingError, req, nil)
				}
				handleError(w, processingError)
				return
			}
		}

		sh.handleResponse(w, req, res)
	}
}

func (sh *SimpleHandler) GetRequestProcessors(path string) []routingI.RequestProcessor {
	rps := make([]routingI.RequestProcessor, 0)
	for _, rp := range sh.requestMiddlewares {
		if rp.GetRegex().MatchString(path) {
			rps = append(rps, rp)
		}
	}
	return rps
}

func (sh *SimpleHandler) GetResponseProcessors(path string) []routingI.ResponseProcessor {
	rps := make([]routingI.ResponseProcessor, 0)
	for _, rp := range sh.responseMiddlewares {
		if rp.GetRegex().MatchString(path) {
			rps = append(rps, rp)
		}
	}
	return rps
}

func (sh *SimpleHandler) AddRequestProcessor(regex string, priority int, processor routingI.RequestProcessor) {
	if sh.requestMiddlewares == nil {
		sh.requestMiddlewares = make([]routingI.RequestProcessor, 0)
	}
	processor.SetRegex(regex)
	processor.SetPriority(priority)
	sh.requestMiddlewares = append(sh.requestMiddlewares, processor)
}

func (sh *SimpleHandler) AddResponseProcessor(regex string, priority int, processor routingI.ResponseProcessor) {
	if sh.responseMiddlewares == nil {
		sh.responseMiddlewares = make([]routingI.ResponseProcessor, 0)
	}
	processor.SetRegex(regex)
	processor.SetPriority(priority)
	sh.responseMiddlewares = append(sh.responseMiddlewares, processor)

}

func (sh *SimpleHandler) GetErrorHandler() *func(error, *routingModel.Request, *routingModel.Response) error {
	return sh.errorHandler
}

func (sh *SimpleHandler) SetErrorHandler(fun *func(error, *routingModel.Request, *routingModel.Response) error) {
	sh.errorHandler = fun
}

func (sh *SimpleHandler) Finalize() {
	logger.Info().Msg("Prioritizing middlewares...")
	//sort middlewares by priority
	if sh.requestMiddlewares != nil {
		sort.Slice(sh.requestMiddlewares, func(i, j int) bool {
			return sh.requestMiddlewares[i].GetPriority() < sh.requestMiddlewares[j].GetPriority()
		})
	}
	if sh.responseMiddlewares != nil {
		sort.Slice(sh.responseMiddlewares, func(i, j int) bool {
			return sh.responseMiddlewares[i].GetPriority() < sh.responseMiddlewares[j].GetPriority()
		})
	}
}
