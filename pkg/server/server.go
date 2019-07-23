package server

import (
	"net/http"

	"gocouchbase/pkg/config"
	"gocouchbase/pkg/log"
	"gocouchbase/pkg/router"
)

type Server struct {
	Config config.Config
	Router router.Router
}

func (s *Server) New() {
	log.Fatal(nil, http.ListenAndServe(
		":"+s.Config.AppPort, s.Router.Mux))
}

func (s *Server) NewTLS() {
	log.Fatal(nil, http.ListenAndServeTLS(
		":"+s.Config.AppPort,
		s.Config.TLSCert,
		s.Config.TLSKey,
		s.Router.Mux))
}
