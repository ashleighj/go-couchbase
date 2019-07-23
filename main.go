package gocouchbase

import (
	"flag"
	"gocouchbase/pkg/api"
	"gocouchbase/pkg/config"
	"gocouchbase/pkg/log"
	"gocouchbase/pkg/router"
	"gocouchbase/pkg/server"
	"gocouchbase/pkg/service"
)

const (
	tlsFlag            = "tls"
	tlsFlagDescription = "Run with TLS server - must have TLS_CERT and TLS_CERT_KEY paths specified in the environment."
)

type app struct {
	server server.Server
}

func main() {
	log.Info(nil, "Hello, server")

	tls := flag.Bool(tlsFlag, false, tlsFlagDescription)
	flag.Parse()

	barebones := app{}

	s := service.NewDefault()
	api.Service = s
	router.Service = s

	r := router.Router{}
	r.Init()

	config := config.Get()
	barebones.server = server.Server{
		Config: config,
		Router: r}

	if *tls {
		barebones.server.NewTLS()
	} else {
		barebones.server.New()
	}
}
