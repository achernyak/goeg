package main

import "encoding/json"

var version int

type JSONInvoice struct {
	Id           int
	CustomerId   int
	DepartmentId string
	Raised       string
	Due          string
	Paid         bool
	Note         string
	Items        []*Item
}

type JSONInvoice100 struct {
	Id         int
	CustomerId int
	Raised     string
	Due        string
	Paid       bool
	Note       string
	Items      []*Item
}

func (invoice Invoice) MarshalJSON() ([]byte, error) {
	jsonInvoice := JSONInvoice{
		invoice.Id,
		invoice.CustomerId,
		invoice.DepartmentId,
		invoice.Raised.Format(dateFormat),
		invoice.Due.Format(dateFormat),
		invoice.Paid,
		invoice.Note,
		invoice.Items,
	}
	return json.Marshal(jsonInvoice)
}
