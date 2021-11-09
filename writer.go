package main

import (
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

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

type color struct {
	r, g, b int
}

func (c color) isWhite() bool {
	return c.r == 255 && c.g == 255 && c.b == 255
}

// setStyleFn returns a func that will write the content with the style passed by the params.
func setStyleFn(style string, size float64, h float64, c color, font string) func(p *gofpdf.Fpdf, c string) {
	return func(pdf *gofpdf.Fpdf, content string) {
		pdf.SetFont(font, style, size)
		pdf.SetFillColor(c.r, c.g, c.b)

		words := strings.Split(content, " ")

		for _, word := range words {
			if len(word) == 0 {
				continue
			}

			finalXPos := pdf.GetStringWidth(" "+word) + pdf.GetX()

			if finalXPos > pageWidth {
				pdf.Ln(h)
			}

			switch style {
			case "B":
				pdf.CellFormat(pdf.GetStringWidth(word)+offsetSpace, h, word, "", 0, "", false, 0, "")
			case "I":
				pdf.CellFormat(pdf.GetStringWidth(word)+offsetSpace, h, word, "", 0, "", false, 0, "")
			default:
				pdf.CellFormat(pdf.GetStringWidth(word)+offsetSpace, h, word, "", 0, "", !c.isWhite(), 0, "")
			}
		}
	}
}

// func that returns the format of the tokens.
func formatWriter(style int) func(p *gofpdf.Fpdf, c string) { // nolint
	switch style {
	case bold:
		return setStyleFn("B", normalTextSize, normalTextHeight, color{255, 255, 255}, "Arial")
	case italic:
		return setStyleFn("I", normalTextSize, normalTextHeight, color{255, 255, 255}, "Arial")
	case code:
		return setStyleFn("", normalTextSize, normalTextHeight, color{220, 220, 220}, "Courier")
	case heading1:
		return setStyleFn("B", heading1Size, headingGrp1Height, color{255, 255, 255}, "Arial")
	case heading2:
		return setStyleFn("B", heading2Size, headingGrp1Height, color{255, 255, 255}, "Arial")
	case heading3:
		return setStyleFn("B", heading3Size, headingGrp2Height, color{255, 255, 255}, "Arial")
	case heading4:
		return setStyleFn("B", heading4Size, headingGrp2Height, color{255, 255, 255}, "Arial")
	case heading5:
		return setStyleFn("B", heading5Size, headingGrp3Height, color{255, 255, 255}, "Arial")
	case heading6:
		return setStyleFn("B", heading6Size, headingGrp3Height, color{255, 255, 255}, "Arial")
	default:
		return setStyleFn("", normalTextSize, normalTextHeight, color{255, 255, 255}, "Arial")
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
			formatWriter(t.style)(p.pdf, t.content)
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
