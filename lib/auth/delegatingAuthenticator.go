package auth

import (
	"errors"
	"strings"

	authModel "github.com/Alquama00s/serverEngine/lib/auth/model"

	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type DelegatingAuthenticator struct {
	authenticators map[string]Authenticator
}

func NewDelegatingAuthenticator() *DelegatingAuthenticator {
	return &DelegatingAuthenticator{
		authenticators: make(map[string]Authenticator),
	}
}

func (d *DelegatingAuthenticator) AddAuthenticator(name string, authenticator Authenticator) *DelegatingAuthenticator {
	d.authenticators[name] = authenticator
	return d
}

func (j *DelegatingAuthenticator) ParsePrincipal(req *routingModel.Request) error {
	if req == nil || req.RawRequest == nil {
		return errors.New("invalid request")
	}

	authHeader := req.RawRequest.Header.Get("Authorization")
	if authHeader == "" {
		req.SetMetaData("auth.principal", authModel.GuestPrincipal())
		return nil
	}

	tokenType := strings.Split(authHeader, " ")
	if len(tokenType) != 2 {
		return errors.New("invalid request")
	}

	auth, exist := j.authenticators[tokenType[0]]
	var err error
	if exist {
		err = auth.ParsePrincipal(req)
	}
	return err
}
