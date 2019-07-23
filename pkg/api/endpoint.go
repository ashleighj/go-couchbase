package api

import (
	"gocouchbase/pkg/log"

	"github.com/go-chi/chi"
)

func RegisterEndpoints(router *chi.Mux) {
	router.Route("/authservice", func(r chi.Router) {

		router.Route("/v1", func(r chi.Router) {

			r.Get("/health", handleHealthCheck)
			r.Post("/authenticate", handleServiceAuthenticate)

			r.Route("/admin", func(r chi.Router) {
				r.Get("/roles", handleRolesGet)
				r.Put("/role", handleRoleCreate)
				r.Post("/client/register", handleRegisterReset)
				r.Post("/client/resetkey", handleRegisterReset)
			})
		})
	})
	log.Infof(nil, log.LogEndpointsRegistered, "API")
}
