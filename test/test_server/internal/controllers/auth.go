package controller

import (
	"net/http"
	"serverEngineTests/internal/model"
	"serverEngineTests/internal/service"

	"github.com/Alquama00s/serverEngine"
	"github.com/Alquama00s/serverEngine/lib/auth"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
	"github.com/Alquama00s/serverEngine/lib/utils/serverUtils"
)

type AuthController struct {
	*service.AuthService
	*auth.JWTAuthenticator
}

func (ac *AuthController) Controllers() {

	ac.JWTAuthenticator = auth.NewJwtAuthenticator()

	router := serverEngine.Registrar().Router("/api/auth")

	serverEngine.Registrar().RegisterSimpleReqProcessor(router.PathPrefix+"/signup",
		func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response) {
			req.Logger.
				Debug().Str("pkg", "controller.AuthController").
				Msg("Unmarshling middleware")
			_, err := serverUtils.Unmarshal[model.User](req)
			if err != nil {
				return nil, err, nil
			}
			return req, nil, nil

		})

	router.Path("/login").
		Method(http.MethodPost).
		Handeler(func(r *routingModel.Request) (*routingModel.Response, error) {
			user, err := serverUtils.Unmarshal[model.User](r)
			if err != nil {
				return nil, err
			}
			u, err := ac.AuthService.Login(user.Username, user.Password)
			if err != nil {
				return nil, err
			}
			role := make([]string, len(u.Roles))
			var priv []string
			for i, r := range u.Roles {
				role[i] = r
				for _, p := range u.Roles {
					priv = append(priv, p)
				}
			}
			tok, err := ac.JWTAuthenticator.CreateToken(priv, role, 1, u.Username)

			if err != nil {
				return nil, err
			}
			return routingModel.NewRestResponse().
				AddToBody("message", "login sueccesfull").
				AddToBody("token", tok), nil
		})

}
