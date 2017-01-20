package omap

import "strings"

type Map struct {
	root   *node
	less   func(interface{}, interface{}) bool
	length int
}

type node struct {
	key, value  interface{}
	red         bool
	left, right *node
}

func New(less func(interface{}, interface{}) bool) *Map {
	return &Map{less: less}
}

func NewCaseFoldedKeyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return strings.ToLower(a.(string)) < strings.ToLower(b.(string))
	}}
}

func NewStringKeyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return a.(string) < b.(string)
	}}
}

func NewIntKeyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}}
}

func NewFloat64Keyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return a.(float64) < b.(float64)
	}}
}
