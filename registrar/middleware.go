package registrar

import (
	"math"
	"net/http"

	"github.com/Alquama00s/serverEngine/lib"
)

func (r *DefaultRegistrar) RegisterRequestProcessors(rps ...lib.RequestProcessor) {
	for _, rp := range rps {
		r.Handler.AddRequestProcessor(rp.GetRegexString(), rp.GetPriority(), rp)
	}
}

func (r *DefaultRegistrar) RegisterResponseProcessors(rps ...lib.ResponseProcessor) {
	for _, rp := range rps {
		r.Handler.AddResponseProcessor(rp.GetRegexString(), rp.GetPriority(), rp)
	}
}

func (r *DefaultRegistrar) RegisterSimpleReqProcessor(regExString string, fun func(req *lib.Request) (*lib.Request, error, *lib.Response)) {
	r.RegisterPrioritizedSimpleReqProcessor(regExString, 0, fun)
}

func (r *DefaultRegistrar) RegisterSimpleResProcessor(regExString string, fun func(res *lib.Response) (*lib.Response, error)) {
	r.RegisterPrioritizedSimpleResProcessor(regExString, 0, fun)
}

func (r *DefaultRegistrar) RegisterPrioritizedSimpleReqProcessor(regExString string, priority int, fun func(req *lib.Request) (*lib.Request, error, *lib.Response)) {
	rp := &lib.SimpleReqMiddleWare{
		RequestProcessor: fun,
	}
	rp.SetRegex(regExString)
	rp.SetPriority(priority)
	r.RegisterRequestProcessors(rp)
}

func (r *DefaultRegistrar) RegisterPrioritizedSimpleResProcessor(regExString string, priority int, fun func(res *lib.Response) (*lib.Response, error)) {
	rp := &lib.SimpleResMiddleWare{
		ResponseProcessor: fun,
	}
	rp.SetRegex(regExString)
	rp.SetPriority(priority)
	r.RegisterResponseProcessors(rp)
}

func (r *DefaultRegistrar) ErrorHandler(fun func(error, *lib.Request, *lib.Response) error) {
	r.Handler.SetErrorHandler(&fun)
}

func (r *DefaultRegistrar) RegisterRouteAuthorizer(routeAuth lib.RouteAuthorizer) {
	r.RegisterRequestProcessors(routeAuth.GetRequestProcessor())
}

func (r *DefaultRegistrar) RegisterAuthenticator(routeAuth lib.Authenticator) {
	r.RegisterPrioritizedSimpleReqProcessor("/*", math.MinInt,
		func(req *lib.Request) (*lib.Request, error, *lib.Response) {
			req.Logger.
				Info().Str("pkg", "registrar.DefaultRegistrar").
				Msg("Authenticating request ..")
			err := routeAuth.ParsePrincipal(req)
			if err != nil {
				_, ok := err.(*lib.ErrorResponse)
				if ok {
					return nil, nil, lib.NewRestResponse().
						SetBody(err).
						SetStatus(http.StatusForbidden)
				}
				return nil, err, lib.NewRestResponse().
					SetBody(lib.NewErrorMessage("Could not authenticate"))
			}

			return req, nil, nil
		})
}
