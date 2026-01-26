package main

import (
	controller "serverEngineTests/internal/controllers"
	"serverEngineTests/internal/initializers"

	"github.com/Alquama00s/serverEngine"
)

func main() {
	//register initializers
	serverEngine.Registrar().
		RegisterInitializers(
			&initializers.CommonMiddleware{},
			&initializers.AuthMiddleware{},
		)
	//register routes
	serverEngine.Registrar().
		RegisterControllerSet(
			&controller.Status{},
			&controller.Options{},
			&controller.AuthController{},
		)
	serverEngine.Port(8888)
	serverEngine.Sereve()
}
