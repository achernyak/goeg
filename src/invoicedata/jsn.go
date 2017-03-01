package main

import (
	"encoding/json"
	"time"
)

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

func (invoice *Invoice) UnmarshalJSON(data []byte) (err error) {
	var jsonInvoice JSONInvoice
	if version == fileVersion {
		if err = json.Unmarshal(data, &jsonInvoice); err != nil {
			return err
		}
	} else {
		var jsonInvoice100 JSONInvoice100
		if err = json.Unmarshal(data, &jsonInvoice100); err != nil {
			return err
		}
		jsonInvoice = JSONInvoice{
			jsonInvoice100.Id,
			jsonInvoice100.CustomerId,
			"",
			jsonInvoice100.Raised,
			jsonInvoice100.Due,
			jsonInvoice100.Paid,
			jsonInvoice100.Note,
			jsonInvoice100.Items,
		}
	}
	var raised, due time.Time
	if raised, err = time.Parse(dateFormat, jsonInvoice.Raised); err != nil {
		return err
	}
	if due, err = time.Parse(dateFormat, jsonInvoice.Due); err != nil {
		return err
	}
	*invoice = Invoice{
		jsonInvoice.Id,
		jsonInvoice.CustomerId,
		jsonInvoice.DepartmentId,
		raised,
		due,
		jsonInvoice.Paid,
		jsonInvoice.Note,
		jsonInvoice.Items,
	}
	return nil
}
