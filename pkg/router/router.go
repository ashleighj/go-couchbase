package router

import (
	"gocouchbase/pkg/api"
	"gocouchbase/pkg/static"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Router struct {
	Mux *chi.Mux
}

func (r *Router) Init() {
	r.initRouter()
	r.initRoutes()
}

func (r *Router) initRouter() {
	r.Mux = chi.NewRouter()

	r.Mux.Use(
		render.SetContentType(render.ContentTypeJSON),
		initContext,
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
}

func (r *Router) initRoutes() {
	api.RegisterEndpoints(r.Mux)
	static.RegisterEndpoints(r.Mux)
}
