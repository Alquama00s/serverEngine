package controller

import (
	"net/http"

	"github.com/Alquama00s/serverEngine"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type Options struct{}

func (a *Options) Controllers() {
	router := serverEngine.Registrar().Router("/")

	router.Path("/").
		Method(http.MethodOptions).
		Handeler(func(r *routingModel.Request) (*routingModel.Response, error) {
			h := &http.Header{}
			return &routingModel.Response{
				StatusCode: 204,
				Headers:    h,
			}, nil
		})

}
