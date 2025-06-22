package lib

type Route struct {
	controller func(*Request) (*Response, error)
	method     string
	path       string
}

func (r *Route) Handeler(handeler func(*Request) (*Response, error)) *Route {
	if handeler == nil {
		panic("handeler cannot be nil")
	}
	if r.controller != nil {
		panic("handeler already set")
	}
	if r.path == "" {
		panic("path cannot be empty")
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
func (r *Route) GetController() func(*Request) (*Response, error) {
	return r.controller
}
