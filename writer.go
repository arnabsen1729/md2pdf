package main

import (
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const pageWidth = 190

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
				pdf.CellFormat(pdf.GetStringWidth(word)+2, h, word, "", 0, "", false, 0, "")
			case "I":
				pdf.CellFormat(pdf.GetStringWidth(word)+2, h, word, "", 0, "", false, 0, "")
			default:
				pdf.CellFormat(pdf.GetStringWidth(word)+2, h, word, "", 0, "", !c.isWhite(), 0, "")
			}
		}
	}
}

// map that contains the format of the tokens.
var formatWriter = map[int]func(p *gofpdf.Fpdf, c string){
	para:     setStyleFn("", 14, 6, color{255, 255, 255}, "Arial"),
	bold:     setStyleFn("B", 14, 6, color{255, 255, 255}, "Arial"),
	italic:   setStyleFn("I", 14, 6, color{255, 255, 255}, "Arial"),
	code:     setStyleFn("", 14, 6, color{220, 220, 220}, "Courier"),
	heading1: setStyleFn("B", 22, 9, color{255, 255, 255}, "Arial"),
	heading2: setStyleFn("B", 20, 9, color{255, 255, 255}, "Arial"),
	heading3: setStyleFn("B", 18, 8, color{255, 255, 255}, "Arial"),
	heading4: setStyleFn("B", 17, 8, color{255, 255, 255}, "Arial"),
	heading5: setStyleFn("B", 16, 7, color{255, 255, 255}, "Arial"),
	heading6: setStyleFn("B", 15, 7, color{255, 255, 255}, "Arial"),
}

type pdfWriter struct {
	pdf *gofpdf.Fpdf
}

func (p *pdfWriter) init(lines [][]*token) {
	p.pdf = gofpdf.New("P", "mm", "A4", "")
	p.pdf.AddPage()

	for _, line := range lines {
		for _, t := range line {
			formatWriter[t.style](p.pdf, t.content)
		}

		p.pdf.Ln(8)
	}
}

func (p *pdfWriter) export(filename string) {
	err := p.pdf.OutputFileAndClose(filename + ".pdf")
	if err != nil {
		log.Fatalln("[ Error occured during exporting pdf ]", err)
	}
}
