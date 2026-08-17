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
	//
	//code, err := qr.Encode("erqwerqw dfg sdfgsdf gsdfgsdfgsdfgsdfgsdfgdf sdfg sdf sdf sdfg sdfgsdfgsdfdf sdfg sdf sdf gsdfgsdfger", qr.Q)
	//if err != nil {
	//	return
	//}
	//
	//data := code.PNG()
	//file, err := os.Create("output.png")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close() // Обязательно закрываем файл в конце
	//
	//// Записываем байты
	//_, err = file.Write(data)
	//if err != nil {
	//	log.Fatal(err)
	//}

}
