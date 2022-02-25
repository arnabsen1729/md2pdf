// Package writer creates a pdf file and takes care of the formatting of the
// tokens.
package writer

import (
	"log"
	"math"
	"os"
	"io"
	"strings"

	"github.com/arnabsen1729/md2pdf/parser"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
)

type color struct {
	r, g, b int
}

// rgb function is used so that the editor can detect color like rgb(0, 0, 255).
func rgb(r int, g int, b int) *color {
	return &color{r, g, b}
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

func displayImage(imageURL string, pdf *gofpdf.Fpdf) {
	// refer to https://stackoverflow.com/questions/51190445/how-to-use-image-url-in-gofpdf.
	pdf.Ln(lineHeight)
	httpimg.Register(pdf, imageURL, "")
	info := pdf.RegisterImage(imageURL, "")
	pdf.Image(imageURL, pdf.GetX(), pdf.GetY(), math.Min(pageWidth, info.Width()), 0, true, "", 0, "")
}

// writeWithStyle will write the content with the style passed by the params.
func writeWithStyle(style string, size float64, h float64, bg *color, font string, fg *color, pdf *gofpdf.Fpdf, t *parser.Token) {
	if t.Style == parser.Image {
		displayImage(t.AltContent, pdf)

		return
	} else if t.Style == parser.Blockquote {
		// prints the small rect in the beginning of blockquote.
		pdf.SetFillColor(200, 200, 200)
		pdf.CellFormat(2, h, " ", "", 0, "", true, 0, "")
	}

	pdf.SetFont(font, style, size)
	pdf.SetTextColor(fg.r, fg.g, fg.b)
	pdf.SetFillColor(bg.r, bg.g, bg.b)

	words := strings.Split(t.Content, " ")

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		finalXPos := pdf.GetStringWidth(" "+word) + pdf.GetX()
		w := pdf.GetStringWidth(word) + offsetSpace

		if finalXPos > pageWidth {
			pdf.Ln(h)

			if t.Style == parser.Blockquote {
				pdf.SetFillColor(200, 200, 200)
				pdf.CellFormat(2, h, " ", "", 0, "", true, 0, "")
				pdf.SetFillColor(bg.r, bg.g, bg.b)
			}
		}

		if t.Style == parser.Link {
			// Currently separate links are added to each word. Ideally the entire phrase should be link.
			pdf.LinkString(pdf.GetX(), pdf.GetY(), w, h, t.AltContent)
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
func formatWriter(p *gofpdf.Fpdf, t *parser.Token) {
	var (
		// standard colors
		black     = rgb(0, 0, 0)
		blue      = rgb(0, 0, 255)
		white     = rgb(255, 255, 255)
		grey      = rgb(220, 220, 220)
		lightGrey = rgb(240, 240, 240)
	)

	switch t.Style {
	case parser.Bold:
		writeWithStyle("B", normalTextSize, normalTextHeight, white, "Helvetica", black, p, t)
	case parser.Italic:
		writeWithStyle("I", normalTextSize, normalTextHeight, white, "Helvetica", black, p, t)
	case parser.Code:
		writeWithStyle("", normalTextSize, normalTextHeight, grey, "Courier", black, p, t)
	case parser.Heading1:
		writeWithStyle("B", heading1Size, headingGrp1Height, white, "Helvetica", black, p, t)
	case parser.Heading2:
		writeWithStyle("B", heading2Size, headingGrp1Height, white, "Helvetica", black, p, t)
	case parser.Heading3:
		writeWithStyle("B", heading3Size, headingGrp2Height, white, "Helvetica", black, p, t)
	case parser.Heading4:
		writeWithStyle("B", heading4Size, headingGrp2Height, white, "Helvetica", black, p, t)
	case parser.Heading5:
		writeWithStyle("B", heading5Size, headingGrp3Height, white, "Helvetica", black, p, t)
	case parser.Heading6:
		writeWithStyle("B", heading6Size, headingGrp3Height, white, "Helvetica", black, p, t)
	case parser.Link:
		writeWithStyle("", normalTextSize, normalTextHeight, white, "Helvetica", blue, p, t)
	case parser.Blockquote:
		writeWithStyle("", normalTextSize, normalTextHeight, lightGrey, "Helvetica", black, p, t)
	default:
		writeWithStyle("", normalTextSize, normalTextHeight, white, "Helvetica", black, p, t)
	}
}

type PdfWriter struct {
	pdf *gofpdf.Fpdf
}

// NewWriter takes the parsed list of list of tokens returns a PdfWriter.
func NewWriter(lines [][]*parser.Token) PdfWriter {
	p := PdfWriter{pdf: nil}
	p.pdf = gofpdf.New("P", "mm", "A4", "")
	p.pdf.AddPage()

	for _, line := range lines {
		for _, t := range line {
			formatWriter(p.pdf, t)
		}

		p.pdf.Ln(lineHeight)
	}

	return p
}

// Export saves the PdfWriter object on disk as a pdf document.
func (p *PdfWriter) Export(filename string) {
	err := p.pdf.OutputFileAndClose(filename + ".pdf")
	if err != nil {
		log.Fatalln("[ Error occurred during exporting pdf ]", err)
		os.Exit(1)
	}
}

// ExportWriter saves the PdfWriter object out to the provided io.Writer.
func (p *PdfWriter) ExportWriter(w io.Writer) {
	err := p.pdf.Output(w)
	if err != nil {
		log.Fatalln("[ Error occurred during exporting pdf ]", err)
		os.Exit(1)
	}
}
