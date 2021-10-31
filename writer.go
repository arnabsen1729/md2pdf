package main

import (
	"log"

	"github.com/jung-kurt/gofpdf"
)

// setStyleFn returns a func that will write the content with the style passed by the params
func setStyleFn(style string, size float64, h float64) func(p *gofpdf.Fpdf, c string) {
	return func(pdf *gofpdf.Fpdf, content string) {
		pdf.SetFont("Arial", style, size)
		pdf.MultiCell(190, h, content, "0", "0", false)
	}
}

// map that contains the format of the tokens
var formatWriter = map[int]func(p *gofpdf.Fpdf, c string){
	para:     setStyleFn("", 14, 6),
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

func (p *pdfWriter) init(tokens []*token) {
	p.pdf = gofpdf.New("P", "mm", "A4", "")
	p.pdf.AddPage()
	for _, t := range tokens {
		formatWriter[t.style](p.pdf, t.content)
	}
}

func (p *pdfWriter) export(filename string) {
	err := p.pdf.OutputFileAndClose(filename + ".pdf")
	if err != nil {
		log.Fatalln("[ Error occured during exporting pdf ]", err)
	}
}
