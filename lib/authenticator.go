package lib

import (
	"encoding/base64"
	"errors"
	"strings"
)

type Authenticator interface {
	ParsePrincipal(*Request) error
}

type BasicAuthenticator struct {
}

func (a *BasicAuthenticator) ParsePrincipal(req *Request) error {
	if req == nil || req.RawRequest == nil {
		return errors.New("invalid request")
	}

	authHeader := req.RawRequest.Header.Get("Authorization")
	if authHeader == "" {
		req.RequestPrincipal = GuestPrincipal()
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
	req.RequestPrincipal = &Principal{
		userName:        username,
		rawToken:        password,
		tokenType:       "Basic",
		isAuthenticated: false,
		privileges:      make(map[string]struct{}),
		roles:           make(map[string]struct{}),
	}

	return nil

}

func NewBasicAuthenticator() *BasicAuthenticator {
	return &BasicAuthenticator{}
}
