package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/arnabsen1729/md2pdf/parser"
	"github.com/arnabsen1729/md2pdf/writer"
)

// readFile to read the contents of the file and return string.
func readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return string(content)
}

func main() {
	fileName := flag.String("file", "", "Name of the markdown file to be converted")
	outputFileName := flag.String("output", "", "Name of the generated PDF file (default: <input-file-name>.pdf)")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *outputFileName == "" {
		*outputFileName = strings.TrimSuffix(*fileName, filepath.Ext(*fileName))
	}

	md := readFile(*fileName)
	par := parser.NewParser(md)
	pdf := writer.NewWriter(par.Lines)
	pdf.Export(*outputFileName)
}
