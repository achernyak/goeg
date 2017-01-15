package main

import "fmt"

type FuzzyBool struct{ value float32 }

func New(value interface{}) (*FuzzyBool, error) {
	amount, err := float32ForValue(value)
	return &FuzzyBool{amount}, err
}

func float32ForValue(value interface{}) (fuzzy float32, err error) {
	switch value := value.(type) {
	case float32:
		fuzzy = value
	case float64:
	case int:
		fuzzy = float32(value)
	case bool:
		fuzzy = 0
		if value {
			fuzzy = 1
		}
	default:
		return 0, fmt.Errorf("float32ForValue(): %v is not a number or boolean",
			value)
	}
	if fuzzy < 0 {
		fuzzy = 0
	} else if fuzzy > 1 {
		fuzzy = 1
	}
	return fuzzy, nil
}

func (fuzzy *FuzzyBool) String() string {
	return fmt.Sprintf("%.0f%%", 100*fuzzy.value)
}

func (fuzzy *FuzzyBool) Set(value interface{}) (err error) {
	fuzzy.value, err = float32ForValue(value)
	return err
}

func (fuzzy *FuzzyBool) Copy() *FuzzyBool {
	return &FuzzyBool{fuzzy.value}
}

func main() {
	a, _ := fuzzybool.New(0)
	b, _ := fuzzybool.New(.25)
	c, _ := fuzzybool.New(.75)
	d := c.Copy()
	if err := d.Set(1); err != nil {
		fmt.Println(err)
	}
	process(a, b, c, d)
	s := []*fuzzybool.FuzzyBool{a, b, c, d}
	fmt.Println(s)
}

func process(a, b, c, d *fuzzybool.FuzzyBool) {
	fmt.Println("Original:", a, b, c, d)
	fmt.Println("Not:     ", a.Not(), b.Not(), c.Not(), d.Not())
	fmt.Println("Not Not: ", a.Not().Not(), b.Not().Not(), c.Not().Not(),
		d.Not().Not())
	fmt.Println("0.And(.25)->", a.And(b), "* .25.And(.75)->", b.And(c),
		"* .75.And(1)->", c.And(d), "* .25.And(.75,1)->", b.And(c, d))
	fmt.Println("0.Or(.25)->", a.Or(b), "* .25.Or(.75)->", b.Or(c),
		"* .75.Or(1)->", c.Or(d), "* .25.Or(.75,1)->", b.Or(c, d))
	fmt.Println("a < c, a == c, a > c:", a.Less(c), a.Equal(c), c.Less(a))
	fmt.Println("Bool:    ", a.Bool(), b.Bool(), c.Bool(), d.Bool())
	fmt.Println("Float:   ", a.Float(), b.Float(), c.Float(), d.Float())
}