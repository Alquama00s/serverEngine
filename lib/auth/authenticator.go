package auth

import (
	"encoding/base64"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/Alquama00s/serverEngine"
	authModel "github.com/Alquama00s/serverEngine/lib/auth/model"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type Authenticator interface {
	ParsePrincipal(*routingModel.Request) error
}

func ApplyAuthenticator(authenticator Authenticator) {
	serverEngine.Registrar().RegisterPrioritizedSimpleReqProcessor("/*", math.MinInt,
		func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response) {
			req.Logger.
				Info().Str("pkg", "registrar.DefaultRegistrar").
				Msg("Authenticating request ..")
			err := authenticator.ParsePrincipal(req)
			if err != nil {
				_, ok := err.(*routingModel.ErrorResponse)
				if ok {
					return nil, nil, routingModel.NewRestResponse().
						SetBody(err).
						SetStatus(http.StatusForbidden)
				}
				return nil, err, routingModel.NewRestResponse().
					SetBody(routingModel.NewErrorMessage("Could not authenticate"))
			}

			return req, nil, nil
		})
}

type BasicAuthenticator struct {
}

func (a *BasicAuthenticator) ParsePrincipal(req *routingModel.Request) error {
	if req == nil || req.RawRequest == nil {
		return errors.New("invalid request")
	}

	authHeader := req.RawRequest.Header.Get("Authorization")
	if authHeader == "" {
		req.SetMetaData("auth.principal", authModel.GuestPrincipal())
		return nil
	}

	// Basic authentication format: "Basic <base64-encoded-credentials>"
	if len(authHeader) < 6 || authHeader[:6] != "Basic " {
		return errors.New("invalid request")
	}

	// Decode the base64 credentials
	credentials, err := base64.StdEncoding.DecodeString(authHeader[6:])
	if err != nil {
		return err
	}

	splittedCreds := strings.SplitN(string(credentials), ":", 2)

	if len(splittedCreds) != 2 {
		return errors.New("invalid request")
	}
	username := strings.TrimSpace(splittedCreds[0])
	password := strings.TrimSpace(splittedCreds[1])
	req.SetMetaData("auth.principal", authModel.NewPrincipal(
		username,
		password,
		"Basic",
		"",
		make(map[string]struct{}),
		make(map[string]struct{}),
		nil,
	))

	return nil

}

func (a *BasicAuthenticator) Apply() {
	ApplyAuthenticator(a)
}

func NewBasicAuthenticator() *BasicAuthenticator {
	return &BasicAuthenticator{}
}
