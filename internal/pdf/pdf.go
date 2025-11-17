package pdf

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"links_project/internal/models"
)

func GeneratePdf(batches []*models.Batch)([]byte, error){
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()
	pdf.Cell(40,10, "Report of links status")
	pdf.Ln(12)

	for _, batch := range batches {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Link_list %d", batch.ID))
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 11)
		for link, status := range batch.Statuses {
			pdf.Cell(0, 8, fmt.Sprintf("%s : %s", link, status))
			pdf.Ln(6)
		}

		pdf.Ln(4)
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}