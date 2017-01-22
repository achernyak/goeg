package oslice

import "strings"

func New(less func(interface{}, interface{}) bool) *Slice {
	return &Slice{less: less}
}

func NewStringSlice() *Slice {
	return &Slice{less: func(a, b interface{}) bool {
		return a.(string) < b.(string)
	}}
}

func NewCaseFoldedSlice() *Slice {
	return &Slice{less: func(a, b interface{}) bool {
		return strings.ToLower(a.(string)) < strings.ToLower(b.(string))
	}}
}

func NewIntSlice() *Slice {
	return &Slice{less: func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}}
}

type Slice struct {
	slice []interface{}
	less  func(interface{}, interface{}) bool
}
