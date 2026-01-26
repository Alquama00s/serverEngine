package controller

import (
	"serverEngineTests/internal/model"

	"github.com/Alquama00s/serverEngine"
	"github.com/Alquama00s/serverEngine/lib/auth"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type Status struct {
}

func (st *Status) Controllers() {
	router := serverEngine.Registrar().Router("/api/status")

	auth.NewSimpleRouteAuth().
		Path(router.PathPrefix).Roles("admin", "alfa").SetAndRoles().
		Apply()

	router.
		Path("/now").Method("GET").
		Handeler(func(r *routingModel.Request) (*routingModel.Response, error) {
			return &routingModel.Response{
				Body: model.NewStatus("hi"),
			}, nil
		})

	router.
		Path("/after").Method("GET").
		Handeler(func(r *routingModel.Request) (*routingModel.Response, error) {
			return &routingModel.Response{
				Body: model.NewStatus("hi after"),
			}, nil
		})
}
