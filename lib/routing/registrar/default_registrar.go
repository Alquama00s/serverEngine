package registrar

import (
	"net/http"
	"strings"

	routingI "github.com/Alquama00s/serverEngine/lib/routing/interface"

	DI "github.com/Alquama00s/serverEngine/lib/DI"
	"github.com/Alquama00s/serverEngine/lib/logging/loggerFactory"
	routing "github.com/Alquama00s/serverEngine/lib/routing"
)

var (
	logger = loggerFactory.GetLogger("DefaultRegistrar")
)

type DefaultRegistrar struct {
	Routers           map[string]*routing.Router
	controllerSetList []routing.ControllerSet
	initializerList   []DI.Initializers
	Handler           routingI.Handler
}

func (r *DefaultRegistrar) Router(prefix string) *routing.Router {
	if val, exist := r.Routers[prefix]; exist {
		return val
	}
	router := &routing.Router{
		PathPrefix: prefix,
		Routes:     []*routing.Route{},
	}
	r.Routers[prefix] = router
	return router
}

func (r *DefaultRegistrar) RegisterControllerSet(cs ...routing.ControllerSet) {
	if r.controllerSetList == nil {
		r.controllerSetList = make([]routing.ControllerSet, 0)
	}
	r.controllerSetList = append(r.controllerSetList, cs...)
}

func (r *DefaultRegistrar) RegisterInitializers(inits ...DI.Initializers) {
	if r.initializerList == nil {
		r.initializerList = make([]DI.Initializers, 0)
	}
	r.initializerList = append(r.initializerList, inits...)
}

func (r *DefaultRegistrar) FinalizeRoutes(mx *http.ServeMux) {
	logger.Info().Msg("Finalizing routes...")
	for _, cs := range r.controllerSetList {
		cs.Controllers()
	}

	for _, router := range r.Routers {
		for _, route := range router.Routes {
			if route.GetController() != nil {
				mx.HandleFunc(route.GetMethod()+" "+strings.TrimSuffix(router.PathPrefix+route.GetPath(), "/"),
					r.Handler.GetHandler(route.GetController()),
				)

			}
		}
	}
	r.Handler.Finalize()
}

func (r *DefaultRegistrar) Initialize() {
	logger.Info().Msg("Initializing...")
	for _, init := range r.initializerList {
		init.Init()
	}
}
