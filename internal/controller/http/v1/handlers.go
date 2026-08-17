package v1

import (
	"crypto/rsa"

	"github.com/mlplabs/lynway-print-service/internal/usecases"

	"github.com/go-chi/chi/v5"
	_ "github.com/mlplabs/common-go-pkg/pkg/http/response/wrapper"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type ServiceHandler struct {
	service *usecases.PrintService
}

func NewServiceHandler(router *chi.Mux, publicKey *rsa.PublicKey, service *usecases.PrintService) {
	//h := &ServiceHandler{
	//	service: service,
	//}
	//wrapper := wrapper.NewWrapper()

	router.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/docs/*", httpSwagger.Handler(
				httpSwagger.URL("doc.json"),
			))
		})
	})
}
