package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
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
