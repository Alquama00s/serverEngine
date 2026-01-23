package routing

import routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"

type Route struct {
	controller func(*routingModel.Request) (*routingModel.Response, error)
	method     string
	path       string
}

func (r *Route) Handeler(handeler func(*routingModel.Request) (*routingModel.Response, error)) *Route {
	if handeler == nil {
		panic("handeler cannot be nil")
	}
	if r.controller != nil {
		panic("handeler already set")
	}
	r.controller = handeler
	return r
}

func (r *Route) Method(method string) *Route {
	if method == "" {
		panic("method cannot be empty")
	}
	if r.method != "" {
		panic("method already set")
	}
	r.method = method
	return r
}

func (r *Route) GetMethod() string {
	return r.method
}

func (r *Route) GetPath() string {
	return r.path
}
func (r *Route) GetController() func(*routingModel.Request) (*routingModel.Response, error) {
	return r.controller
}
