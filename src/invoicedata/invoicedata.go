package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	magicNumber = 0x125D
	fileVersion = 101
	fileType    = "INVOICES"
	dateFormat  = "2006-01-02" // This date must always be used (see text).

)

type Invoice struct {
	Id           int
	CustomerId   int
	DepartmentId string
	Raised       time.Time
	Due          time.Time
	Paid         bool
	Note         string
	Item         []*Item
}

type Item struct {
	Id       string
	Price    float64
	Quantity int
	TaxBand  int
	Note     string
}

type InvoicesMarshaler interface {
	MarshalInvoices(writer io.Writer, invoices []*Invoice) error
}

type InvoicesUnmarshaler interface {
	UnmarshalInvoices(reader io.Reader) ([]*Invoice, error)
}

func main() {
	log.SetFlags(0)
	if len(os.Args) != 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		log.Fatalf("usage: %s infile.ext outfile.ext\n"+
			".ext may be .json or .txt optionally gzipped (e.g., .gob.gz)\n", filepath.Base(os.Args[0]))
	}
	inFilename, outFilename := os.Args[1], os.Args[2]
	if inFilename == outFilename {
		log.Fatalln("won't overwrite a file with itself")
	}

	invoices, err := readInvoiceFile(inFilename)
	if err != nil {
		log.Fatalln("Faild to read:", err)
	}
	if err := writeInvoiceFile(outFilename, invoices); err != nil {
		log.Fatalln("Failed to write:", err)
	}
}

func readInvoiceFile(filename string) ([]*Invoice, error) {
	file, closer, err := openInvoiceFile(filename)
	if closer != nil {
		defer closer()
	}
	if err != nil {
		return nil, err
	}
	return readInvoices(file, suffixOf(filename))
}

func openInvoiceFile(filename string) (io.ReadCloser, func(), error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	closer := func() { file.Close() }
	var reader io.ReadCloser = file
	var decompressor *gzip.Reader
	if strings.HasSuffix(filename, ".gz") {
		if decompressor, err = gzip.NewReader(file); err != nil {
			return file, closer, err
		}
		closer = func() { decompressor.Close(); file.Close() }
		reader = decompressor
	}
	return reader, closer, nil
}

func readInvoices(reader io.Reader, suffix string) ([]*Invoice, error) {
	var unmarshaler InvoicesUnmarshaler
	switch suffix {
	case ".jsn", ".json":
		unmarshaler = JSONMarshaler{}
	case ".txt":
		unmarshaler = TxtMarshaler{}
	}
	if unmarshaler != nil {
		return unmarshaler.UnmarshalInvoices(reader)
	}
	return nil, fmt.Errorf("unrecognized input suffix: %s", suffix)
}

func writeInvoices(writer io.Writer, suffix string,
	invoices []*Invoice) error {
	var marshaler InvoicesMarshaler
	switch suffix {
	case ".jsn", ".json":
		marshaler = JSONMarshaler{}
	case ".txt":
		marshaler = TxtMarshaler{}
	}
	if marshaler != nil {
		return marshaler.MarshalInvoices(writer, invoices)
	}
	return errors.New("unrecognized output suffix")
}
