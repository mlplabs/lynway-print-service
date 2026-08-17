package ipp

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/mlplabs/lynway-print-service/internal/model"

	"github.com/hiventive/go-ipp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/charmap"
	_ "rsc.io/qr/coding"
)

type Client struct {
	ippClient  *ipp.IPPClient
	cupsClient *ipp.CUPSClient
}

func NewClient(cfg *Config) *Client {
	// 1. Создаем IPP-клиент.
	// CUPS по умолчанию работает на localhost:631
	clientIpp := ipp.NewIPPClient(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.UseTLS)
	clientCups := ipp.NewCUPSClient(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.UseTLS)
	return &Client{ippClient: clientIpp, cupsClient: clientCups}
}

// GetPrinters получение списка принтеров из CUPS
func (c *Client) GetPrinters() ([]Printer, error) {
	printers, err := c.cupsClient.GetPrinters(nil)
	if err != nil {
		return nil, fmt.Errorf("error getting printers: %w", err)
	}
	res := make([]Printer, 0)
	for name, _ := range printers {
		res = append(res, Printer{Name: name})
	}
	return res, nil
}

// PrintLabelViaIPP печать стандартной этикетки на принтер
func (c *Client) PrintLabelViaIPP(pdfData *bytes.Buffer, printer string, format model.LabelFormat) error {
	var docName, mimeType, media string
	switch format {
	case model.LabelFormat_43_25:
		docName = "Label_43x25_Job"
		mimeType = "application/pdf"
		media = "Standard.43x25mm"
	default:
		docName = "Label_Unknown_Job"
		mimeType = "application/pdf"
		media = "Unknown"
	}
	doc := ipp.Document{
		Document: pdfData,
		Size:     pdfData.Len(),
		Name:     docName,
		MimeType: mimeType,
	}

	options := map[string]interface{}{
		"media": media,
	}

	_, _, err := c.ippClient.PrintJob(doc, printer, options)
	if err != nil {
		return fmt.Errorf("ошибка выполнения PrintJob: %w", err)
	}

	return nil
}

// --- TEST ---

func (c *Client) PrintImg(printerName string) error {
	// 1. Создаем графический холст (43x25 мм = 344x200 пикселей при 203 DPI)
	imgWidth := 344
	imgHeight := 200
	canvas := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Заливаем фон чистым белым цветом
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// Отрисовываем текст и графику средствами Go
	addText(canvas, 20, 40, "Go-IPP PrintFile")
	addText(canvas, 20, 80, "Format: image/png")
	addText(canvas, 20, 120, "Universal OK")
	drawCircle(canvas, 270, 60, 25)

	// 2. Кодируем изображение в формат PNG
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, canvas); err != nil {
		log.Fatalf("Ошибка кодирования в PNG: %v", err)
	}

	// 3. Сохраняем во временный файл (так как метод PrintFile работает с диском)
	tempFileName := "temp_label.png"
	err := os.WriteFile(tempFileName, pngBuf.Bytes(), 0666)
	if err != nil {
		log.Fatalf("Ошибка записи временного файла: %v", err)
	}
	defer os.Remove(tempFileName) // Удаляем файл после завершения программы

	// 4. Инициализируем стандартный IPP клиент phin1x
	client := ipp.NewIPPClient("localhost", 631, "", "", false)

	// 🔥 ИСПРАВЛЕНИЕ ТИПА: Используем строго map[string]interface{}
	jobAttributes := map[string]interface{}{
		"document-format": "image/png",
	}

	// 5. Отправляем файл на печать через высокоуровневый метод PrintFile
	// 🔥 ИСПРАВЛЕНИЕ: Метод возвращает (int, error), где первое число — это напрямую Job ID
	fmt.Printf("Отправка PNG файла на принтер '%s' через IPP.PrintFile...\n", printerName)
	jobID, err := client.PrintFile(tempFileName, printerName, jobAttributes)
	if err != nil {
		log.Fatalf("Критическая ошибка IPP: %v", err)
	}

	// 6. Выводим ID созданной задачи в CUPS
	fmt.Printf("Успешно! Задача добавлена в очередь CUPS. Job ID: %d\n", jobID)

	return nil
}

func drawCircle(img *image.RGBA, cx, cy, r int) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r {
				img.Set(x, y, color.Black)
			}
		}
	}
}

func addText(img *image.RGBA, x, y int, label string) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: inconsolata.Bold8x16,
		// Функция fixed.P(x, y) создает правильную структуру fixed.Point26_6
		Dot: fixed.P(x, y),
	}
	d.DrawString(label)
}

// sudo chmod 666 /dev/usb/lp1 - разово
//sudo usermod -aG lp имя_пользователя

func (c *Client) PrintTSPL(printerName string) error {

	utf8Str := "Привет из Go!"

	// Конвертируем её в Windows-1251
	encoder := charmap.Windows1251.NewEncoder()
	win1251Bytes, err := encoder.Bytes([]byte(utf8Str))
	if err != nil {
		log.Fatal(err)
	}

	tsplScript := "SIZE 43 mm, 25 mm\r\n" +
		"GAP 2 mm, 0 mm\r\n" +
		"DIRECTION 1\r\n" +
		"CLS\r\n" +
		"CODEPAGE 1251\r\n" +
		"TEXT 50,50,\"3\",0,1,1,\"" + string(win1251Bytes) + "\"\r\n" +
		"BARCODE 50,160,\"128\",50,1,0,2,2,\"12345678\"\r\n" +
		//"SOUND 2, 200\r\n" + // Заставит принтер дважды пискнуть
		"PRINT 1, 1\r\n"

	f, err := os.OpenFile("/dev/usb/lp1", os.O_WRONLY, 0)
	if err != nil {
		log.Fatalf("Не удалось открыть принтер: %v", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(tsplScript))
	if err != nil {
		log.Fatalf("Ошибка записи: %v", err)
	}

	_ = f.Sync()
	return nil
}
