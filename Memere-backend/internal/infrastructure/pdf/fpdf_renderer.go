// Package pdf provides PDF generation infrastructure for course certificates.
// The Renderer interface is defined in the certificate usecase package so that
// clean-architecture dependency rules are preserved (usecases never import this
// package; only the composition-root / main.go does).
package pdf

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-pdf/fpdf"

	certificate "github.com/Bereke1t2/Memere/memere-backend/internal/usecase/certificate"
)

// FPDFRenderer generates course completion certificates using the pure-Go fpdf
// library. No headless browser or external binary dependency.
type FPDFRenderer struct{}

var _ certificate.Renderer = FPDFRenderer{}

// NewFPDFRenderer returns a ready-to-use renderer.
func NewFPDFRenderer() FPDFRenderer { return FPDFRenderer{} }

// Certificate renders a completion certificate PDF and returns the raw bytes.
func (FPDFRenderer) Certificate(_ context.Context, data certificate.CertData) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "A4", "") // landscape A4
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// --- Background border ---
	pdf.SetDrawColor(0, 100, 0)
	pdf.SetLineWidth(3)
	pdf.Rect(10, 10, 277, 190, "D")

	// --- Logo / header ---
	pdf.SetFont("Arial", "B", 28)
	pdf.SetTextColor(0, 100, 0)
	pdf.CellFormat(277, 20, "Memere — ExamPrep", "", 1, "C", false, 0, "")

	// --- Title ---
	pdf.SetFont("Arial", "", 16)
	pdf.SetTextColor(80, 80, 80)
	pdf.Ln(6)
	pdf.CellFormat(277, 10, "Certificate of Completion", "", 1, "C", false, 0, "")

	// --- Body ---
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(40, 40, 40)
	pdf.CellFormat(277, 10, "This certifies that", "", 1, "C", false, 0, "")

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(277, 14, data.StudentName, "", 1, "C", false, 0, "")

	pdf.Ln(4)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(40, 40, 40)
	pdf.CellFormat(277, 10, "has successfully completed", "", 1, "C", false, 0, "")

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 60, 120)
	pdf.CellFormat(277, 12, data.CourseTitle, "", 1, "C", false, 0, "")

	// --- Date ---
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(277, 8, fmt.Sprintf("Issued on %s", data.IssuedDate), "", 1, "C", false, 0, "")

	// --- Serial (verification code) ---
	pdf.Ln(6)
	pdf.SetFont("Courier", "", 10)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(277, 8, fmt.Sprintf("Verification code: %s", data.Serial), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf: render certificate: %w", err)
	}
	return buf.Bytes(), nil
}
