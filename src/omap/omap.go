package omap

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
