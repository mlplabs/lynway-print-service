package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mlplabs/lynway-print-service/internal/config"
	controller "github.com/mlplabs/lynway-print-service/internal/controller/http"
	v1 "github.com/mlplabs/lynway-print-service/internal/controller/http/v1"
	"github.com/mlplabs/lynway-print-service/internal/di"
	"github.com/mlplabs/lynway-print-service/internal/usecases"
	"github.com/mlplabs/lynway-print-service/internal/workers/get_print_queue"

	"github.com/go-chi/chi/v5"
	httpServer "github.com/mlplabs/common-go-pkg/pkg/http/server"
)

func Run(cfg *config.Config) {

	ctx := context.Background()
	di := di.NewDi(cfg)
	service := usecases.NewPrintService(di.GetIPPClient())

	get_print_queue.NewWorker(&cfg.GetPrintQueueWorker, service, di.GetWhsClient()).Do(ctx)

	router := controller.NewRouter()
	apiRouter := chi.NewRouter()

	v1.NewServiceHandler(apiRouter, nil, service)

	router.Mount("/api", apiRouter)

	server := httpServer.NewServer(&router, httpServer.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)

	log.Println("app.Run: service successful started")

	select {
	case s := <-interrupt:
		log.Printf("app.Run: terminated, signal: %s", s.String())
	case err := <-server.Notify():
		log.Printf("app.Run: httpServer.Notify: %s", err)

		return
	}

	err := server.Shutdown()
	if err != nil {
		log.Printf("app.Run: httpServer.Shutdown: %s", err)
		return
	}
}
