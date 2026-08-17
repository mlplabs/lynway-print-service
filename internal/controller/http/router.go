package http

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// NewRouter возвращает новый роутер.
func NewRouter() chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.NoCache,
		cors.AllowAll().Handler,
		Language,
	)

	// config CORS for allowing all origins
	router.Use(
		middleware.SetHeader("Access-Control-Allow-Origin", "*"),
		middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, HEAD, PUT"),
		middleware.SetHeader("Access-Control-Allow-Headers",
			"Access-Control-Allow-Headers, Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Content-Length, Accept-Encoding, Authorization, ResponseType"),
		middleware.SetHeader("Access-Control-Allow-Credentials", "true"),
		middleware.SetHeader("Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type"),
		middleware.SetHeader("Access-Control-Allow-Max-Age", "86400"),
	)

	router.Get("/health", health)

	return *router
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
