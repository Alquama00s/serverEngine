package registrar

import (
	routingI "github.com/Alquama00s/serverEngine/lib/routing/interface"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"

	routing "github.com/Alquama00s/serverEngine/lib/routing"
)

func (r *DefaultRegistrar) RegisterRequestProcessors(rps ...routingI.RequestProcessor) {
	for _, rp := range rps {
		r.Handler.AddRequestProcessor(rp.GetRegexString(), rp.GetPriority(), rp)
	}
}

func (r *DefaultRegistrar) RegisterResponseProcessors(rps ...routingI.ResponseProcessor) {
	for _, rp := range rps {
		r.Handler.AddResponseProcessor(rp.GetRegexString(), rp.GetPriority(), rp)
	}
}

func (r *DefaultRegistrar) RegisterSimpleReqProcessor(regExString string, fun func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response)) {
	r.RegisterPrioritizedSimpleReqProcessor(regExString, 0, fun)
}

func (r *DefaultRegistrar) RegisterSimpleResProcessor(regExString string, fun func(res *routingModel.Response) (*routingModel.Response, error)) {
	r.RegisterPrioritizedSimpleResProcessor(regExString, 0, fun)
}

func (r *DefaultRegistrar) RegisterPrioritizedSimpleReqProcessor(regExString string, priority int, fun func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response)) {
	rp := &routing.SimpleReqMiddleWare{
		RequestProcessor: fun,
	}
	rp.SetRegex(regExString)
	rp.SetPriority(priority)
	r.RegisterRequestProcessors(rp)
}

func (r *DefaultRegistrar) RegisterPrioritizedSimpleResProcessor(regExString string, priority int, fun func(res *routingModel.Response) (*routingModel.Response, error)) {
	rp := &routing.SimpleResMiddleWare{
		ResponseProcessor: fun,
	}
	rp.SetRegex(regExString)
	rp.SetPriority(priority)
	r.RegisterResponseProcessors(rp)
}

func (r *DefaultRegistrar) ErrorHandler(fun func(error, *routingModel.Request, *routingModel.Response) error) {
	r.Handler.SetErrorHandler(&fun)
}

// func (r *DefaultRegistrar) RegisterRouteAuthorizer(routeAuth routingModel.RouteAuthorizer) {
// 	r.RegisterRequestProcessors(routeAuth.GetRequestProcessor())
// }

// func (r *DefaultRegistrar) RegisterAuthenticator(routeAuth routingModel.Authenticator) {
// 	r.RegisterPrioritizedSimpleReqProcessor("/*", math.MinInt,
// 		func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response) {
// 			req.Logger.
// 				Info().Str("pkg", "registrar.DefaultRegistrar").
// 				Msg("Authenticating request ..")
// 			err := routeAuth.ParsePrincipal(req)
// 			if err != nil {
// 				_, ok := err.(*routingModel.ErrorResponse)
// 				if ok {
// 					return nil, nil, routingModel.NewRestResponse().
// 						SetBody(err).
// 						SetStatus(http.StatusForbidden)
// 				}
// 				return nil, err, routingModel.NewRestResponse().
// 					SetBody(routingModel.NewErrorMessage("Could not authenticate"))
// 			}

// 			return req, nil, nil
// 		})
// }
