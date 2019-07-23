package auth

import (
	"gocouchbase/pkg/config"
	"time"

	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(config.Get().JWTSecret), nil)
}

func Get() *jwtauth.JWTAuth {
	return tokenAuth
}

func GetEncodedToken(clientId string, expiration time.Time) (tokenString string, err error) {
	_, tokenString, err = tokenAuth.Encode(jwtauth.Claims{
		config.ClientIDKey: clientId,
		"exp":              expiration.Unix(),
		"iat":              time.Now().Unix()})

	return
}
