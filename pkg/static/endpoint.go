package static

import (
	"net/http"
	"os"
	"path/filepath"

	"gocouchbase/pkg/log"

	"github.com/go-chi/chi"
)

func RegisterEndpoints(router *chi.Mux) {
	workDir, _ := os.Getwd()
	staticFilesDir := http.Dir(filepath.Join(workDir, "pkg/static"))

	router.Handle("/", http.FileServer(staticFilesDir))

	log.Infof(nil, log.LogEndpointsRegistered, "Static")
}
