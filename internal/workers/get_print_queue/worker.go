package get_print_queue

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mlplabs/lynway-print-service/internal/model"
	"github.com/mlplabs/lynway-print-service/internal/usecases"

	whs "github.com/mlplabs/lynway-wms/pkg/client"
	"github.com/mlplabs/lynway-wms/pkg/dto"
)

type GetPrintQueueWorker struct {
	cfg       *Config
	service   *usecases.PrintService
	whsClient *whs.HttpClient
}

func NewWorker(
	cfg *Config,
	service *usecases.PrintService,
	whsClient *whs.HttpClient,
) *GetPrintQueueWorker {
	return &GetPrintQueueWorker{
		cfg:       cfg,
		service:   service,
		whsClient: whsClient,
	}
}

func (w *GetPrintQueueWorker) Do(ctx context.Context) {
	if !w.cfg.Enabled {
		log.Printf("DEBUG: %s worker is disabled", w.cfg.Name)

		return
	}

	w.process(ctx)

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				log.Printf("ERROR: %s worker. context done, err %v", w.cfg.Name, err)
			}

			return
		case <-time.After(w.cfg.Interval):
			w.process(ctx)
		}
	}
}

func (w *GetPrintQueueWorker) process(ctx context.Context) {
	log.Println("DEBUG: GetPrintQueueWorker process start")
	//w.whsClient.UpdateAgentInfo()
	//
	printers, err := w.service.GetPrinters()
	if err != nil {
		log.Printf("ERROR: GetPrintQueueWorker process get printers error, err %v", err)
	}
	info := dto.AgentInfoRequest{}
	for _, printer := range printers {
		info.Printers = append(info.Printers, dto.PrinterInfo{
			PrinterName: printer.Name,
		})
	}

	err = w.whsClient.UpdateAgentInfo(ctx, &info)
	if err != nil {
		log.Printf("ERROR: GetPrintQueueWorker process update agent info error, err %v", err)
	}

	d, err := w.whsClient.GetPrintTask(ctx, nil)
	if err != nil {
		log.Printf("ERROR: %s get print task err %v", w.cfg.Name, err)
		return
	}
	buf, err := w.service.GenerateLabelPdf(ctx, &d.Data, d.Labels, model.LabelFormat_43_25)
	if err != nil {
		log.Printf("ERROR: %s generate label pdf err %v", w.cfg.Name, err)
	}
	if d.ToFile {
		err = os.WriteFile("output.pdf", buf.Bytes(), 0644)
		if err != nil {
			fmt.Printf("ERROR: write output pdf err: %v", err)
		}
	} else {
		err = w.service.PrintPdfLabel(buf, d.PrinterName, model.LabelFormat_43_25)
		if err != nil {
			log.Printf("ERROR: %s print label err %v", w.cfg.Name, err)
		}
	}

	if err != nil {
		w.whsClient.UpdatePrintTaskStatus(ctx, &dto.UpdatePrintTaskStatusRequest{
			Id:     d.Id,
			Status: dto.PrintTaskStatusSuccess,
		})

		return
	}

	w.whsClient.UpdatePrintTaskStatus(ctx, &dto.UpdatePrintTaskStatusRequest{
		Id:     d.Id,
		Status: dto.PrintTaskStatusSuccess,
	})
}
