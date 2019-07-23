package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gocouchbase/pkg/api"
	"gocouchbase/pkg/config"
	"gocouchbase/pkg/log"
	"gocouchbase/pkg/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
)

var Service service.Service

func initContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		id := uuid.New()

		rCtx := map[string]string{
			config.RequestIDKey:   id.String(),
			config.RequestTimeKey: fmt.Sprint(time.Now().UTC()),
			config.RemoteAddrKey:  r.RemoteAddr,
			config.MethodKey:      r.Method,
			config.RequestURIKey:  r.RequestURI}

		ctx = context.WithValue(ctx, config.RequestContext, rCtx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func validate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		path := r.URL.Path

		for _, p := range config.OpenPaths {
			if strings.Contains(path, p) {
				log.Info(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		if !strings.Contains(path, config.APIPrefix) {
			api.RespondWithError(w, r, http.StatusBadRequest, log.ErrMissingRQPrefix)
			return
		}

		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			api.RespondWithError(w, r, http.StatusUnauthorized, log.ErrTokenExpired)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			api.RespondWithError(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		log.Info(ctx, claims)

		client, err := Service.GetAPIClient(claims[config.ClientIDKey].(string))
		if err != nil {
			api.RespondWithError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		if client.ClientID == "" {
			api.RespondWithError(w, r, http.StatusUnauthorized, log.ErrEndpointAuth)
			return
		}

		role, err := Service.GetServiceRole(client.Role)
		if role.RoleName == "" {
			api.RespondWithError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		for _, rolePath := range role.AllowedPaths {
			matched := match(rolePath, path)

			if matched {
				next.ServeHTTP(w, r)
				return
			}
		}

		api.RespondWithError(w, r, http.StatusUnauthorized, log.ErrEndpointAuth)
		return
	}
	return http.HandlerFunc(fn)
}

func match(pattern string, path string) bool {
	path = strings.TrimPrefix(path, config.APIPrefix)
	wildcard := "*"

	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	for i, patternPart := range patternParts {
		if patternPart != wildcard && patternPart != pathParts[i] {
			return false
		}
		if patternPart == wildcard {
			return true
		}
	}
	return false
}
