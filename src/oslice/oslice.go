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

func (slice *Slice) Clear() {
	slice.slice = nil
}

func (slice *Slice) Add(x interface{}) {
	if slice.slice == nil {
		slice.slice = []interface{}{x}
	} else if index := bisectLeft(slice.slice, slice.less, x); index == len(slice.slice) {
		slice.slice = append(slice.slice, x)
	} else {
		updateSlice := make([]interface{}, len(slice.slice)+1)
		at := copy(updateSlice, slice.slice[:index])
		at += copy(updateSlice[at:], []interface{}{x})
		copy(updateSlice[at:], slice.slice[index:])
		slice.slice = updateSlice
	}
}

func (slice *Slice) Remove(x interface{}) bool {
	index := bisectLeft(slice.slice, slice.less, x)
	for ; index < len(slice.slice); index++ {
		if !slice.less(slice.slice[index], x) &&
			!slice.less(x, slice.slice[index]) {
			slice.slice = append(slice.slice[:index],
				slice.slice[index+1:]...)
			return true
		}
	}
	return false
}

func (slice *Slice) Index(x interface{}) int {
	index := bisectLeft(slice.slice, slice.less, x)
	if index >= len(slice.slice) ||
		slice.less(slice.slice[index], x) ||
		slice.less(x, slice.slice[index]) {
		return -1
	}
	return index
}

func (slice *Slice) At(index int) interface{} {
	return slice.slice[index]
}

func (slice *Slice) Len() int {
	return len(slice.slice)
}

func bisectLeft(slice []interface{},
	less func(interface{}, interface{}) bool, x interface{}) int {
	left, right := 0, len(slice)
	for left < right {
		middle := int((left + right) / 2)
		if less(slice[middle], x) {
			left = middle + 1
		} else {
			right = middle
		}
	}
	return left
}
