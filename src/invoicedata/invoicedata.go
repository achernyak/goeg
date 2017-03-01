package main

import (
	"fmt"
	"io"
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
	fmt.Println("vim-go")
}
