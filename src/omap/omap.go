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

func (m *Map) Insert(key, value interface{}) (inserted bool) {
	m.root, inserted = m.insert(m.root, key, value)
	m.root.red = false
	if inserted {
		m.length++
	}
	return inserted
}

func (m *Map) insert(root *node, key, value interface{}) (*node, bool) {
	inserted := false
	if root == nil {
		return &node{key: key, value: value, red: true}, true
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	if m.less(key, root.key) {
		root.left, inserted = m.insert(root.left, key, value)
	} else if m.less(root.key, key) {
		root.right, inserted = m.insert(root.right, key, value)
	} else {
		root.value = value
	}
	if isRed(root.right) && !isRed(root.left) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	return root, inserted
}

func isRed(root *node) bool { return root != nil && root.red }

func colorFlip(root *node) {
	root.red = !root.red
	if root.left != nil {
		root.left.red = !root.left.red
	}
	if root.right != nil {
		root.right.red = !root.right.red
	}
}

func rotateLeft(root *node) *node {
	x := root.right
	root.right = x.left
	x.left = root
	x.red = root.red
	root.red = true
	return x
}

func rotateRight(root *node) *node {
	x := root.left
	root.left = x.right
	x.right = root
	x.red = root.red
	root.red = true
	return x
type Option struct {
	Fill   color.Color
type Option struct {
	Fill   color.Color
	Radius int
}

func New(shape string, option Option) (Shaper, error) {
	sidesForShape := map[string]int{"triangle": 3, "square": 4,
		"pentagon": 5, "hexagon": 6, "heptagon": 7, "octagon": 8,
		"enneagon": 9, "nonagon": 9, "decagon": 10}
	if sides, found := sidesForShape[shape]; found {
		return NewRegularPolygon(option.Fill, option.Radius, sides), nil
	}
	if shape != "circle" {
		return nil, fmt.Errorf("shapes.New(): invalide shape '%s'", shape)
	}
	return NewCircle(option.Fill, option.Radius), nil
}

	Radius int
}

func New(shape string, option Option) (Shaper, error) {
	sidesForShape := map[string]int{"triangle": 3, "square": 4,
		"pentagon": 5, "hexagon": 6, "heptagon": 7, "octagon": 8,
		"enneagon": 9, "nonagon": 9, "decagon": 10}
	if sides, found := sidesForShape[shape]; found {
		return NewRegularPolygon(option.Fill, option.Radius, sides), nil
	}
	if shape != "circle" {
		return nil, fmt.Errorf("shapes.New(): invalide shape '%s'", shape)
	}
	return NewCircle(option.Fill, option.Radius), nil
}

}
