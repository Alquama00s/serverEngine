package lib

import (
	"errors"
	"strings"
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

func (j *DelegatingAuthenticator) ParsePrincipal(req *Request) error {
	if req == nil || req.RawRequest == nil {
		return errors.New("invalid request")
	}

	authHeader := req.RawRequest.Header.Get("Authorization")
	if authHeader == "" {
		req.RequestPrincipal = GuestPrincipal()
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
