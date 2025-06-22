package serverEngine

import (
	"log"
	"net/http"
	"sync"

	"github.com/Alquama00s/serverEngine/lib"
	"github.com/Alquama00s/serverEngine/loggerFactory"
	"github.com/Alquama00s/serverEngine/registrar"
)

var (
	routeRegistrar *registrar.DefaultRegistrar
	once           sync.Once
	logger         = loggerFactory.GetLogger()
)

func Registrar() *registrar.DefaultRegistrar {
	once.Do(func() {
		logger.Info().Msg("Initializing route registrar")
		if routeRegistrar == nil {
			routeRegistrar = &registrar.DefaultRegistrar{
				Routers: make(map[string]*lib.Router),
				Handler: &lib.SimpleHandler{},
			}
		}
	})
	return routeRegistrar
}

func Sereve() {
	logger.Info().Msg("Starting server...")
	routeRegistrar.Initialize()
	routeRegistrar.FinalizeRoutes()
	logger.Info().Msg("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
