package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
)

type JWTAuthenticator struct {
	webKeys             []*jose.JSONWebKey
	signatureAlgorithms []jose.SignatureAlgorithm
	tokenExpiryDuration time.Duration
}

type JWTClaims struct {
	Roles     []string
	Privilege []string
	UserId    uint
	jwt.Claims
}

var (
	newJWTAuthenticator *JWTAuthenticator
	jwtOnce             sync.Once
)

const (
	INVALID_CREDS = "invalid credentials or credential expired"
)

func NewJwtAuthenticator() *JWTAuthenticator {
	jwtOnce.Do(func() {
		sigAlgo := []jose.SignatureAlgorithm{jose.RS256, jose.RS384}
		keys := make([]*jose.JSONWebKey, len(sigAlgo))

		for i := range sigAlgo {
			rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				panic(err)
			}
			keys[i] = &jose.JSONWebKey{
				Key:       rsaKey,
				KeyID:     uuid.NewString() + strconv.Itoa(i),
				Algorithm: string(sigAlgo[i]),
				Use:       "sig",
			}
		}

		newJWTAuthenticator = &JWTAuthenticator{
			webKeys:             keys,
			signatureAlgorithms: sigAlgo,
			tokenExpiryDuration: time.Hour * 24,
		}
	})
	return newJWTAuthenticator
}

func (j *JWTAuthenticator) ParsePrincipal(req *Request) error {
	if req == nil || req.RawRequest == nil {
		return NewErrorMessage(INVALID_CREDS)
	}

	authHeader := req.RawRequest.Header.Get("Authorization")
	if authHeader == "" {
		req.RequestPrincipal = GuestPrincipal()
		return nil
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return NewErrorMessage(INVALID_CREDS)
	}

	parsedJwt, err := jwt.ParseSigned(authHeader[7:], j.signatureAlgorithms)
	if err != nil {
		return NewErrorMessage(INVALID_CREDS)
	}
	if len(parsedJwt.Headers[0].KeyID) < 36 {
		return NewErrorMessage(INVALID_CREDS)
	}
	kIndex, err := strconv.Atoi(string(parsedJwt.Headers[0].KeyID[36]))
	if err != nil {
		return err
	}
	var claim JWTClaims
	err = parsedJwt.
		Claims(&(j.webKeys[kIndex].Key.(*rsa.PrivateKey).PublicKey), &claim)

	if err != nil {
		return NewErrorMessage(INVALID_CREDS)
	}

	priv := make(map[string]struct{})
	role := make(map[string]struct{})

	for _, r := range claim.Roles {
		role[r] = struct{}{}

	}
	for _, p := range claim.Privilege {
		priv[p] = struct{}{}
	}

	req.RequestPrincipal = NewAuthenticatedPrincipal(claim.Subject, authHeader[7:],
		"Bearer", priv, role, claim)

	return nil
}

func (j *JWTAuthenticator) CreateToken(priv, role []string, userId uint, userName string) (string, error) {
	// role := make([]string, len(user.Roles))
	// var priv []string
	// for i, r := range user.Roles {
	// 	role[i] = r.Name
	// 	for _, p := range r.Privileges {
	// 		priv = append(priv, p.Name)
	// 	}
	// }

	claims := JWTClaims{
		Roles:     role,
		Privilege: priv,
		UserId:    userId,
		Claims: jwt.Claims{
			Issuer:    "assist",
			Audience:  []string{userName},
			Subject:   userName,
			Expiry:    jwt.NewNumericDate(time.Now().Add(j.tokenExpiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(j.webKeys))))
	if err != nil {
		return "", err
	}
	key := j.webKeys[i.Int64()]

	signer, err := jose.NewSigner(jose.SigningKey{Key: key, Algorithm: jose.SignatureAlgorithm(key.Algorithm)}, nil)
	if err != nil {
		return "", err
	}

	token, err := jwt.Signed(signer).Claims(claims).Serialize()

	return token, err
}
