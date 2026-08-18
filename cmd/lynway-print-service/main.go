package main

import (
	"log"

	"github.com/mlplabs/lynway-print-service/internal/app"
	"github.com/mlplabs/lynway-print-service/internal/config"

	"github.com/joho/godotenv"
	_ "rsc.io/qr/coding"
)

func main() {
	godotenv.Load(".env")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ошибка чтения конфигурации: %v", err)
	}
	//cfg.AppVersion = AppVersion
	app.Run(cfg)
}
