package di

import (
	"github.com/mlplabs/lynway-print-service/internal/adapters/clients/ipp"
	"github.com/mlplabs/lynway-print-service/internal/config"

	whs "github.com/mlplabs/lynway-wms/pkg/client"
)

type Di struct {
	cfg       *config.Config
	ippClient *ipp.Client
	whsClient *whs.HttpClient
}

func NewDi(cfg *config.Config) *Di {
	return &Di{cfg: cfg}
}

func (d *Di) GetIPPClient() *ipp.Client {
	if d.ippClient != nil {
		return d.ippClient
	}
	d.ippClient = ipp.NewClient(&d.cfg.Ipp)
	return d.ippClient
}

func (d *Di) GetWhsClient() *whs.HttpClient {
	if d.whsClient != nil {
		return d.whsClient
	}
	d.whsClient = whs.NewHttpClient("", &d.cfg.Whs)
	return d.whsClient
}
