package lib

import (
	"fmt"
	"math"
	"net/http"
)

type RouteAuthorizer interface {
	GetRequestProcessor() RequestProcessor
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

func (sr *SimpleRouteAuthorizer) GetRequestProcessor() RequestProcessor {
	rp := &SimpleReqMiddleWare{}
	rp.SetPriority(math.MinInt + 1)
	rp.SetRegex(sr.pathRegex)
	rp.Process(func(req *Request) (*Request, error, *Response) {

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

		if req.RequestPrincipal == nil || !req.RequestPrincipal.IsAuthenticated() {
			return req, nil, NewRestResponse().
				SetBody(NewErrorMessage("Unauthorized User is not authenticated")).
				SetStatus(http.StatusUnauthorized)
		}

		if len(sr.privileges) > 0 {
			if sr.isPrevilegeAnd {
				if !req.RequestPrincipal.AndPrivilegeAuth(sr.privileges...) {
					return nil, nil, NewRestResponse().SetBody(NewErrorMessage("does not have all required privileges")).
						SetStatus(http.StatusUnauthorized)
				}
			} else {
				if !req.RequestPrincipal.OrPrivilegeAuth(sr.privileges...) {
					return nil, nil, NewRestResponse().SetBody(NewErrorMessage("does not have required privileges")).
						SetStatus(http.StatusUnauthorized)
				}
			}
		}

		if len(sr.roles) > 0 {
			if sr.isRoleAnd {
				if !req.RequestPrincipal.AndRoleAuth(sr.roles...) {
					return nil, nil, NewRestResponse().SetBody(NewErrorMessage("does not have all required roles")).
						SetStatus(http.StatusUnauthorized)
				}
			} else {
				if !req.RequestPrincipal.OrRoleAuth(sr.roles...) {
					return nil, nil, NewRestResponse().SetBody(NewErrorMessage("does not have required roles")).
						SetStatus(http.StatusUnauthorized)
				}
			}
		}

		if !sr.IsTokenValid(req.RequestPrincipal.tokenType) {
			return nil, nil, NewRestResponse().
				SetBody(NewErrorMessage(fmt.Sprintf("requre one of %s authentication", sr.__tokenTypes))).
				SetStatus(http.StatusUnauthorized)
		}

		return req, nil, nil
	})

	return rp
}
