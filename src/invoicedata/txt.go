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
