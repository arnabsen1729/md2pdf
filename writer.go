package main

import (
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type color struct {
	r, g, b int
}

func (c color) isWhite() bool {
	return c.r == 255 && c.g == 255 && c.b == 255
}

const (
	pageWidth         = 190
	offsetSpace       = 2
	lineHeight        = 8
	normalTextSize    = 14
	heading1Size      = 22
	heading2Size      = 20
	heading3Size      = 18
	heading4Size      = 17
	heading5Size      = 16
	heading6Size      = 15
	normalTextHeight  = 6
	headingGrp1Height = 9
	headingGrp2Height = 8
	headingGrp3Height = 7
)

// setStyleFn returns a func that will write the content with the style passed by the params.
func setStyleFn(style string, size float64, h float64, bg color, font string, fg color, pdf *gofpdf.Fpdf, t *token) {
	pdf.SetFont(font, style, size)
	pdf.SetTextColor(fg.r, fg.g, fg.b)
	pdf.SetFillColor(bg.r, bg.g, bg.b)

	words := strings.Split(t.content, " ")

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		finalXPos := pdf.GetStringWidth(" "+word) + pdf.GetX()
		w := pdf.GetStringWidth(word) + offsetSpace

		if finalXPos > pageWidth {
			pdf.Ln(h)
		}

		if t.style == link {
			// Currently separate links are added to each word. Ideally the entire phrase should be link.
			pdf.LinkString(pdf.GetX(), pdf.GetY(), w, h, t.altContent)
		}

		switch style {
		case "B":
			pdf.CellFormat(w, h, word, "", 0, "", false, 0, "")
		case "I":
			pdf.CellFormat(w, h, word, "", 0, "", false, 0, "")
		default:
			pdf.CellFormat(w, h, word, "", 0, "", !bg.isWhite(), 0, "")
		}
	}
}

// func that returns the format of the tokens.
func formatWriter(p *gofpdf.Fpdf, t *token) { // nolint
	var (
		// standard colors
		black = color{0, 0, 0}
		blue  = color{0, 0, 255}
		white = color{255, 255, 255}
	)

	switch t.style {
	case bold:
		setStyleFn("B", normalTextSize, normalTextHeight, white, "Arial", black, p, t)
	case italic:
		setStyleFn("I", normalTextSize, normalTextHeight, white, "Arial", black, p, t)
	case code:
		setStyleFn("", normalTextSize, normalTextHeight, color{220, 220, 220}, "Courier", black, p, t)
	case heading1:
		setStyleFn("B", heading1Size, headingGrp1Height, white, "Arial", black, p, t)
	case heading2:
		setStyleFn("B", heading2Size, headingGrp1Height, white, "Arial", black, p, t)
	case heading3:
		setStyleFn("B", heading3Size, headingGrp2Height, white, "Arial", black, p, t)
	case heading4:
		setStyleFn("B", heading4Size, headingGrp2Height, white, "Arial", black, p, t)
	case heading5:
		setStyleFn("B", heading5Size, headingGrp3Height, white, "Arial", black, p, t)
	case heading6:
		setStyleFn("B", heading6Size, headingGrp3Height, white, "Arial", black, p, t)
	case link:
		setStyleFn("", normalTextSize, normalTextHeight, white, "Arial", blue, p, t)
	default:
		setStyleFn("", normalTextSize, normalTextHeight, white, "Arial", black, p, t)
	}
}

type pdfWriter struct {
	pdf *gofpdf.Fpdf
}

func (p *pdfWriter) init(lines [][]*token) {
	p.pdf = gofpdf.New("P", "mm", "A4", "")
	p.pdf.AddPage()

	for _, line := range lines {
		for _, t := range line {
			formatWriter(p.pdf, t)
		}

		p.pdf.Ln(lineHeight)
	}
}

func (p *pdfWriter) export(filename string) {
	err := p.pdf.OutputFileAndClose(filename + ".pdf")
	if err != nil {
		log.Fatalln("[ Error occurred during exporting pdf ]", err)
	}
}
