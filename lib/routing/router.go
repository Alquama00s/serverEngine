package routing

type Router struct {
	PathPrefix string
	Routes     []*Route
}

func (r *Router) Path(path string) *Route {
	route := &Route{
		path:       path,
		controller: nil,
	}
	r.Routes = append(r.Routes, route)
	return route
}
