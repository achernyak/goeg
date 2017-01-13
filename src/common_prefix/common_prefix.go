package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	testData := [][]string{
		{"/home/user/goeg", "/home/user/goeg/prefix",
			"/home/user/goeg/prefix/extra"},
		{"/home/user/goeg", "/home/user/goeg/prefix",
			"/home/user/prefix/extra"},
		{"/pecan/π/goeg", "/pecan/π/goeg/prefix",
			"/pecan/π/prefix/extra"},
		{"/pecan/π/circle", "/pecan/π/circle/prefix",
			"/pecan/π/circle/prefix/extra"},
		{"/home/user/goeg", "/home/users/goeg",
			"/home/userspace/goeg"},
		{"/home/user/goeg", "/tmp/user", "/var/log"},
		{"/home/mark/goeg", "/home/user/goeg"},
		{"home/user/goeg", "/tmp/user", "/var/log"},
	}
	for _, data := range testData {
		fmt.Printf("[")
		gap := ""
		for _, datum := range data {
			fmt.Printf("%s\"%s\"", gap, datum)
			gap = " "
		}
		fmt.Println("]")
		cp := CommonPrefix(data)
		cpp := CommonPathPrefix(data)
		equal := "=="
		if cpp != cp {
			equal = "!="
		}
		fmt.Printf("char X path prefix: \"%s\" %s \"%s\"\n\n",
			cp, equal, cpp)
	}
}

func CommonPrefix(texts []string) string {
	components := make([][]rune, len(texts))
	for i, text := range texts {
		components[i] = []rune(text)
	}
	if len(components) == 0 || len(components[0]) == 0 {
		return ""
	}
	var common bytes.Buffer
FINISH:
	for column := 0; column < len(components[0]); column++ {
		char := components[0][column]
		for row := 1; row < len(components); row++ {
			if column >= len(components[row]) ||
				components[row][column] != char {
				break FINISH
			}
		}
		common.WriteRune(char)
	}
	return common.String()
}

func CommonPathPrefix(paths []string) string {
	const separator = string(filepath.Separator)
	components := make([][]string, len(paths))
	for i, path := range paths {
		components[i] = strings.Split(path, separator)
		if strings.HasPrefix(path, separator) {
			components[i] = append([]string{separator}, components[i]...)
		}
	}
	if len(components) == 0 || len(components[0]) == 0 {
		return ""
	}
	var common []string
FINISH:
	for column := range components[0] {
		part := components[0][column]
		for row := 1; row < len(components); row++ {
			if len(components[row]) == 0 ||
				column >= len(components[row]) ||
				components[row][column] != part {
				break FINISH
			}
		}
		common = append(common, part)
	}
	return filepath.Join(common...)
}
