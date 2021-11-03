package main

import (
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const pageWidth = 190

// setStyleFn returns a func that will write the content with the style passed by the params.
func setStyleFn(style string, size float64, h float64) func(p *gofpdf.Fpdf, c string, i bool) {
	return func(pdf *gofpdf.Fpdf, content string, isStart bool) {
		pdf.SetFont("Arial", style, size)

		words := strings.Split(content, " ")

		for index, word := range words {
			finalXPos := pdf.GetStringWidth(" "+word) + pdf.GetX()

			if finalXPos > pageWidth {
				pdf.Ln(h)
			} else if !(isStart && index == 0) {
				word = " " + word
			}

			pdf.Cell(pdf.GetStringWidth(word), h, word)
		}
	}
}

// map that contains the format of the tokens.
var formatWriter = map[int]func(p *gofpdf.Fpdf, c string, i bool){
	para:     setStyleFn("", 14, 6),
	bold:     setStyleFn("B", 14, 6),
	italic:   setStyleFn("I", 14, 6),
	heading1: setStyleFn("B", 22, 9),
	heading2: setStyleFn("B", 20, 9),
	heading3: setStyleFn("B", 18, 8),
	heading4: setStyleFn("B", 17, 8),
	heading5: setStyleFn("B", 16, 7),
	heading6: setStyleFn("B", 15, 7),
}

type pdfWriter struct {
	pdf *gofpdf.Fpdf
}

func (p *pdfWriter) init(lines [][]*token) {
	p.pdf = gofpdf.New("P", "mm", "A4", "")
	p.pdf.AddPage()

	for _, line := range lines {
		for index, t := range line {
			formatWriter[t.style](p.pdf, t.content, index == 0)
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
