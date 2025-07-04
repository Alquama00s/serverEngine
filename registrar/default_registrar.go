package registrar

import (
	"net/http"
	"strings"

	"github.com/Alquama00s/serverEngine/lib"
	"github.com/Alquama00s/serverEngine/loggerFactory"
)

var (
	logger = loggerFactory.GetLogger("DefaultRegistrar")
)

type DefaultRegistrar struct {
	Routers           map[string]*lib.Router
	controllerSetList []lib.ControllerSet
	initializerList   []lib.Initializers
	Handler           lib.Handler
}

func (r *DefaultRegistrar) Router(prefix string) *lib.Router {
	if val, exist := r.Routers[prefix]; exist {
		return val
	}
	router := &lib.Router{
		PathPrefix: prefix,
		Routes:     []*lib.Route{},
	}
	r.Routers[prefix] = router
	return router
}

func (r *DefaultRegistrar) RegisterControllerSet(cs ...lib.ControllerSet) {
	if r.controllerSetList == nil {
		r.controllerSetList = make([]lib.ControllerSet, 0)
	}
	r.controllerSetList = append(r.controllerSetList, cs...)
}

func (r *DefaultRegistrar) RegisterInitializers(inits ...lib.Initializers) {
	if r.initializerList == nil {
		r.initializerList = make([]lib.Initializers, 0)
	}
	r.initializerList = append(r.initializerList, inits...)
}

func (r *DefaultRegistrar) FinalizeRoutes() {
	logger.Info().Msg("Finalizing routes...")
	for _, cs := range r.controllerSetList {
		cs.Controllers()
	}

	for _, router := range r.Routers {
		for _, route := range router.Routes {
			if route.GetController() != nil {
				http.HandleFunc(route.GetMethod()+" "+strings.TrimSuffix(router.PathPrefix+route.GetPath(), "/"),
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
