package usecases

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mlplabs/lynway-print-service/internal/adapters/clients/ipp"
	"github.com/mlplabs/lynway-print-service/internal/model"
	pdffont "github.com/mlplabs/lynway-print-service/internal/usecases/font"

	"github.com/go-pdf/fpdf"
	"rsc.io/qr"
)

type PrintService struct {
	ippClient *ipp.Client
}

func NewPrintService(ippClient *ipp.Client) *PrintService {
	return &PrintService{
		ippClient: ippClient,
	}
}

func (s *PrintService) GetPrinters() ([]ipp.Printer, error) {
	return s.ippClient.GetPrinters()
}

func (s *PrintService) PrintPdfLabel(data *bytes.Buffer, printerName string, format model.LabelFormat) error {
	return s.ippClient.PrintLabelViaIPP(data, printerName, format)
}

func (s *PrintService) GenerateLabelPdf(ctx context.Context, qrData *string, labels map[string]string, format model.LabelFormat) (*bytes.Buffer, error) {
	switch format {
	case model.LabelFormat_43_25:
		return s.generateLabelPdf4325(qrData, labels)
	}
	return nil, nil
}

// generateLabelPdf4325 формирование pdf
// TODO: QR код планируется сделать опциональным
func (s *PrintService) generateLabelPdf4325(qrData *string, labels map[string]string) (*bytes.Buffer, error) {
	var labelWidth = 43.0  // Ширина этикетки в мм
	var labelHeight = 25.0 // Высота этикетки в мм

	initType := fpdf.InitType{
		UnitStr: "mm",
		Size:    fpdf.SizeType{Wd: labelWidth, Ht: labelHeight},
	}

	regular, _ := base64.StdEncoding.DecodeString(pdffont.ArialBase64)
	bold, _ := base64.StdEncoding.DecodeString(pdffont.ArialBoldBase64)

	pdf := fpdf.NewCustom(&initType)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddUTF8FontFromBytes("Arial", "", regular)
	pdf.AddUTF8FontFromBytes("Arial", "B", bold)
	// Устанавливаем минимальные отступы для маленькой этикетки
	//pdf.SetMargins(1.5, 1.5, 1.5)
	pdf.AddPage()

	qrText := ""
	if qrData != nil {
		qrText = *qrData
	}

	code, err := qr.Encode(qrText, qr.Q)
	if err != nil {
		return nil, err
	}

	qrReader := bytes.NewReader(code.PNG())

	// Настраиваем параметры изображения. В ImageType передаем формат (PNG или JPG)
	imageOptions := fpdf.ImageOptions{
		ImageType: "PNG",
	}

	// Регистрируем картинку под произвольным текстовым ключом, например "my_qr_code"
	pdf.RegisterImageOptionsReader("my_qr_code", imageOptions, qrReader)
	if pdf.Err() {
		return nil, fmt.Errorf("ошибка регистрации изображения в PDF: %v", pdf.Err())
	}

	// Позиционируем QR-код на маленькой этикетке справа
	qrSize := 25.0
	qrX := labelWidth - qrSize
	qrY := (labelHeight - qrSize) / 2.0 // Центрируем по вертикали

	// Выводим изображение на страницу
	pdf.ImageOptions("my_qr_code", qrX, qrY, qrSize, qrSize, false, imageOptions, 0, "")
	// --------------------------------------------

	// Выводим текстовую информацию слева от QR-кода
	textWidth := qrX - 3.0

	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(1.5, 3.0)
	pdf.CellFormat(textWidth, 5.5, labels["title"], "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7)
	pdf.SetX(1.5)
	pdf.CellFormat(textWidth, 3.5, labels["subtitle"], "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 7)
	pdf.SetX(1.5)
	pdf.CellFormat(textWidth, 3.5, labels["description"], "", 1, "L", false, 0, "")

	pdf.SetX(1.5)
	pdf.CellFormat(textWidth, 3.5, labels["comment"], "", 1, "L", false, 0, "")

	// Записываем результат работы fpdf в bytes.Buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
