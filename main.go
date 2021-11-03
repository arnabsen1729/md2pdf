package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// readFile to read the contents of the file and return string
func readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return string(content)
}

func main() {
	fileName := flag.String("file", "", "Name of the markdown file to read")
	outputFileName := flag.String("output", "", "Name of the PDF file to be exported  (default: <input-file-name>.pdf)")
	flag.Parse()

	if *fileName == "" {
		log.Fatal("Please provide a file to read")
		os.Exit(1)
	}

	if *outputFileName == "" {
		*outputFileName = strings.TrimSuffix(*fileName, filepath.Ext(*fileName))
	}

	md := readFile(*fileName)
	p := newParser(md)
	pdf := &pdfWriter{}
	pdf.init(p.lines)
	pdf.export(*outputFileName)
}
