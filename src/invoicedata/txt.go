package main

import (
	"bufio"
	"fmt"
	"io"
)

const noteSep = ":"

type TxtMarshaler struct{}

func (TxtMarshaler) MarshalInvoices(writer io.Writer,
	invoices []*Invoice) error {
	bufferedWriter := bufio.NewWriter(writer)
	defer bufferedWriter.Flush()
	var write writerFunc = func(format string,
		args ...interface{}) error {
		_, err := fmt.Fprintf(bufferedWriter, format, args...)
		return err
	}
	if err := write("%s %d\n", fileType, fileVersion); err != nil {
		return err
	}
	for _, invoice := range invoices {
		if err := write.writeInvoice(invoice); err != nil {
			return err
		}
	}
	return nil
}

type writerFunc func(string, ...interface{}) error

func (write writerFunc) writeInvoice(invoice *Invoice) error {
	note := ""
	if invoice.Note != "" {
		note = noteSep + " " + invoice.Note
	}
	if err := write("INVOICE ID=%d CUSTOMER=%d DEPARTMENT=%s RAISED=%s "+
		"DUE=%s PAID=%t%s\n", invoice.Id, invoice.CustomerId,
		invoice.DepartmentId, invoice.Raised.Format(dateFormat),
		invoice.Due.Format(dateFormat), invoice.Paid, note); err != nil {
		return err
	}
	if err := write.writeItems(invoice.Items); err != nil {
		return err
	}
	return write("\f\n")
}

func (write writerFunc) writeItems(items []*Item) error {
	for _, item := range items {
		if item.Note != "" {
			note = noteSep + " " + item.Note
		}
		if err := write("ITEM ID=%s PRICE=%.2f QUANTITY=%d TAXBAND=%d%s\n",
			item.Id, item.Price, item.Quantity, item.TaxBand,
			note); err != nil {
			return err
		}
	}
	return nil
}

func (TxtMarshaler) UnmarshalInvoices(reader io.Reader) (
	invoices []*Invoice, err error) {
	bufferedReader := bufio.NewReader(reader)
	var version int
	if version, err = checkTxtVersion(bufferedReader); err != nil {
		return nil, err
	}
	var line string
	eof := false
	for lino := 2; !eof; lino++ {
		line, err = bufferedReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return nil, err
		}
		if invoices, err = parseTxtLine(version, lino, line,
			invoices); err != nil {
			return nil, err
		}
	}
	return invoices, nil
}
