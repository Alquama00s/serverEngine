package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/Alquama00s/serverEngine/loggerFactory"
)

var (
	logger = loggerFactory.GetLogger("SimpleHandler")
)

type SimpleHandler struct {
	requestMiddlewares  []RequestProcessor
	responseMiddlewares []ResponseProcessor
	errorHandler        *func(error, *Request, *Response) error
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	byt, err := json.Marshal(NewError(err))
	if err != nil {
		fmt.Printf("error: %v", err)
		w.Write([]byte("error marshling"))
	} else {
		w.Write(byt)
	}
}

func (sh *SimpleHandler) handleResponse(w http.ResponseWriter, req *Request, res *Response) {
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

func (sh *SimpleHandler) GetHandler(controller func(*Request) (*Response, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lgr := loggerFactory.GetLogger(r.URL.Path, r.Host)
		req := &Request{
			RawRequest: r,
			Logger:     lgr,
		}
		var processingError error
		var res *Response

		var tempRes *Response
		var tempReq *Request

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

func (sh *SimpleHandler) GetRequestProcessors(path string) []RequestProcessor {
	rps := make([]RequestProcessor, 0)
	for _, rp := range sh.requestMiddlewares {
		if rp.GetRegex().MatchString(path) {
			rps = append(rps, rp)
		}
	}
	return rps
}

func (sh *SimpleHandler) GetResponseProcessors(path string) []ResponseProcessor {
	rps := make([]ResponseProcessor, 0)
	for _, rp := range sh.responseMiddlewares {
		if rp.GetRegex().MatchString(path) {
			rps = append(rps, rp)
		}
	}
	return rps
}

func (sh *SimpleHandler) AddRequestProcessor(regex string, priority int, processor RequestProcessor) {
	if sh.requestMiddlewares == nil {
		sh.requestMiddlewares = make([]RequestProcessor, 0)
	}
	processor.SetRegex(regex)
	processor.SetPriority(priority)
	sh.requestMiddlewares = append(sh.requestMiddlewares, processor)
}

func (sh *SimpleHandler) AddResponseProcessor(regex string, priority int, processor ResponseProcessor) {
	if sh.responseMiddlewares == nil {
		sh.responseMiddlewares = make([]ResponseProcessor, 0)
	}
	processor.SetRegex(regex)
	processor.SetPriority(priority)
	sh.responseMiddlewares = append(sh.responseMiddlewares, processor)

}

func (sh *SimpleHandler) GetErrorHandler() *func(error, *Request, *Response) error {
	return sh.errorHandler
}

func (sh *SimpleHandler) SetErrorHandler(fun *func(error, *Request, *Response) error) {
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
