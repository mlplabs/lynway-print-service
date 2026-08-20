package config

import (
	"github.com/mlplabs/lynway-print-service/internal/adapters/clients/ipp"
	"github.com/mlplabs/lynway-print-service/internal/workers/get_print_queue"

	whs "github.com/mlplabs/lynway-wms/pkg/client"
)

type (
	Config struct {
		AppVersion          string
		HTTP                `envPrefix:"HTTP_"`
		Whs                 whs.Config             `envPrefix:"WHS_"`
		Ipp                 ipp.Config             `envPrefix:"IPP_"`
		GetPrintQueueWorker get_print_queue.Config `envPrefix:"GET_PRINT_QUEUE_"`
	}
	HTTP struct {
		Host string `env:"HOST" envDefault:"localhost"`
		Port string `env:"PORT" envDefault:"9321"`
	}
)
