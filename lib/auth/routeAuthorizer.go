package auth

import (
	"fmt"
	"math"
	"net/http"

	authModel "github.com/Alquama00s/serverEngine/lib/auth/model"
	routing "github.com/Alquama00s/serverEngine/lib/routing"
	routingI "github.com/Alquama00s/serverEngine/lib/routing/interface"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type RouteAuthorizer interface {
	GetRequestProcessor() routingI.RequestProcessor
}

type SimpleRouteAuthorizer struct {
	privileges     []string
	roles          []string
	tokenTypes     map[string]struct{}
	__tokenTypes   []string
	pathRegex      string
	isPrevilegeAnd bool
	isRoleAnd      bool
}

func NewSimpleRouteAuth() *SimpleRouteAuthorizer {
	return &SimpleRouteAuthorizer{}
}

func (sr *SimpleRouteAuthorizer) Path(pathRegex string) *SimpleRouteAuthorizer {
	sr.pathRegex = pathRegex
	return sr
}

func (sr *SimpleRouteAuthorizer) Privileges(privileges ...string) *SimpleRouteAuthorizer {
	sr.privileges = privileges
	return sr
}

func (sr *SimpleRouteAuthorizer) TokenType(tokenTypes ...string) *SimpleRouteAuthorizer {
	if sr.tokenTypes == nil {
		sr.tokenTypes = make(map[string]struct{})
	}
	for _, tt := range tokenTypes {
		sr.tokenTypes[tt] = struct{}{}
		sr.__tokenTypes = append(sr.__tokenTypes, tt)
	}
	return sr
}

func (sr *SimpleRouteAuthorizer) Roles(roles ...string) *SimpleRouteAuthorizer {
	sr.roles = roles
	return sr
}

func (sr *SimpleRouteAuthorizer) SetOrRoles() *SimpleRouteAuthorizer {
	sr.isRoleAnd = false
	return sr
}

func (sr *SimpleRouteAuthorizer) SetAndRoles() *SimpleRouteAuthorizer {
	sr.isRoleAnd = true
	return sr
}

func (sr *SimpleRouteAuthorizer) SetOrPrivilege() *SimpleRouteAuthorizer {
	sr.isPrevilegeAnd = false
	return sr
}

func (sr *SimpleRouteAuthorizer) SetAndPrivilege() *SimpleRouteAuthorizer {
	sr.isPrevilegeAnd = true
	return sr
}

func (sr *SimpleRouteAuthorizer) IsTokenValid(token string) bool {
	if sr.tokenTypes != nil && len(sr.tokenTypes) >= 0 {
		_, exist := sr.tokenTypes[token]
		return exist
	}
	return true
}

func (sr *SimpleRouteAuthorizer) GetRequestProcessor() routingI.RequestProcessor {
	rp := &routing.SimpleReqMiddleWare{}
	rp.SetPriority(math.MinInt + 1)
	rp.SetRegex(sr.pathRegex)
	rp.Process(func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response) {

		req.Logger.Info().
			Str("pkg", "serverEngine.lib.SimpleRouteAuthorizer").
			Msg("Performing Authorization")
		req.Logger.Info().
			Str("pkg", "serverEngine.lib.SimpleRouteAuthorizer").
			Msg(req.RawRequest.Method)
		if req.RawRequest.Method == http.MethodOptions {
			req.Logger.Info().
				Str("pkg", "serverEngine.lib.SimpleRouteAuthorizer").
				Msg("omitting for options")
			return req, nil, nil
		}

		principal := authModel.ParsePrincipal(req)

		if principal == nil || !principal.IsAuthenticated() {
			return req, nil, routingModel.NewRestResponse().
				SetBody(routingModel.NewErrorMessage("Unauthorized User is not authenticated")).
				SetStatus(http.StatusUnauthorized)
		}

		if len(sr.privileges) > 0 {
			if sr.isPrevilegeAnd {
				if !principal.AndPrivilegeAuth(sr.privileges...) {
					return nil, nil, routingModel.NewRestResponse().SetBody(routingModel.NewErrorMessage("does not have all required privileges")).
						SetStatus(http.StatusUnauthorized)
				}
			} else {
				if !principal.OrPrivilegeAuth(sr.privileges...) {
					return nil, nil, routingModel.NewRestResponse().SetBody(routingModel.NewErrorMessage("does not have required privileges")).
						SetStatus(http.StatusUnauthorized)
				}
			}
		}

		if len(sr.roles) > 0 {
			if sr.isRoleAnd {
				if !principal.AndRoleAuth(sr.roles...) {
					return nil, nil, routingModel.NewRestResponse().SetBody(routingModel.NewErrorMessage("does not have all required roles")).
						SetStatus(http.StatusUnauthorized)
				}
			} else {
				if !principal.OrRoleAuth(sr.roles...) {
					return nil, nil, routingModel.NewRestResponse().SetBody(routingModel.NewErrorMessage("does not have required roles")).
						SetStatus(http.StatusUnauthorized)
				}
			}
		}

		if !sr.IsTokenValid(principal.GetTokenType()) {
			return nil, nil, routingModel.NewRestResponse().
				SetBody(routingModel.NewErrorMessage(fmt.Sprintf("requre one of %s authentication", sr.__tokenTypes))).
				SetStatus(http.StatusUnauthorized)
		}

		return req, nil, nil
	})

	return rp
}
