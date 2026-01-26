package initializers

import (
	"math"
	"serverEngineTests/internal/constants"
	"serverEngineTests/internal/service"

	"github.com/Alquama00s/serverEngine"
	"github.com/Alquama00s/serverEngine/lib/auth"
	authModel "github.com/Alquama00s/serverEngine/lib/auth/model"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

// @Init
type AuthMiddleware struct {
}

func (m *AuthMiddleware) Init() {
	authService := service.NewAuthService()

	//registor authenticators
	auth.NewJwtAuthenticator().Apply()
	// serverEngine.Registrar().RegisterAuthenticator(lib.NewDelegatingAuthenticator().
	// 	AddAuthenticator("Basic", lib.NewBasicAuthenticator()).
	// 	AddAuthenticator("Bearer", lib.NewJwtAuthenticator()))

	//basic authentication
	serverEngine.Registrar().RegisterPrioritizedSimpleReqProcessor("/*", math.MinInt+1, func(r *routingModel.Request) (*routingModel.Request, error, *routingModel.Response) {
		if r == nil {
			r.SetMetaData(constants.AUTH_STATUS, constants.UNAUTHENTICATED)
			return r, nil, nil
		}
		principal := authModel.ParsePrincipal(r)
		if principal == nil {
			r.SetMetaData(constants.AUTH_STATUS, constants.UNAUTHENTICATED)
			return r, nil, nil
		}

		r.Logger.
			Debug().Str("pkg", "controller.GlobalConfigs").
			Msg(principal.GetUserName() + " is trying to authenticate")
		if principal == authModel.GuestPrincipal() || principal.GetTokenType() != "Basic" {
			return r, nil, nil
		}

		u, err := authService.Login(principal.GetUserName(), principal.GetToken())
		if err != nil {
			r.SetMetaData(constants.AUTH_STATUS, constants.UNAUTHENTICATED)
			return r, nil, nil
		}

		priv := make(map[string]struct{})
		role := make(map[string]struct{})

		for _, r := range u.Roles {
			role[r] = struct{}{}
			for _, p := range u.Roles {
				priv[p] = struct{}{}
			}
		}

		principal.Authenticate()
		principal.Apply(r)
		// r.RequestPrincipal = lib.NewAuthenticatedPrincipal(
		// 	u.Username, r.RequestPrincipal.GetToken(),
		// 	r.RequestPrincipal.GetTokenType(), u.ID, priv, role, nil,
		// )
		return r, nil, nil
	})

	// service.NewAccessService().InitializeAcesses()
}
