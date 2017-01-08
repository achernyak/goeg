package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

func main() {
	irregularMatrix := [][]int{{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11},
		{12, 13, 14, 15},
		{16, 17, 18, 19, 20}}
	fmt.Println("irregular:", irregularMatrix)
	slice := Flatten(irregularMatrix)
	fmt.Printf("1x%d: %v\n", len(slice), slice)
	fmt.Printf(" 3x%d: %v\n", neededRows(slice, 3), Make2D(slice, 3))
	fmt.Printf(" 4x%d: %v\n", neededRows(slice, 4), Make2D(slice, 4))
	fmt.Printf(" 5x%d: %v\n", neededRows(slice, 5), Make2D(slice, 5))
	fmt.Printf(" 6x%d: %v\n", neededRows(slice, 6), Make2D(slice, 6))
	slice = []int{9, 1, 9, 5, 4, 4, 2, 1, 5, 4, 8, 8, 4, 3, 6, 9, 5, 7, 5}
	fmt.Println("Original:", slice)
	slice = UniqueInts(slice)
	fmt.Println("Unique:  ", slice)

	iniData := []string{
		"; Cut down copy of Mozilla application.ini file",
		"",
		"[App]",
		"Vendor=Mozilla",
		"Name=Iceweasel",
		"Profile=mozilla/firefox",
		"Version=3.5.16",
		"[Gecko]",
		"MinVersion=1.9.1",
		"MaxVersion=1.9.1.*",
		"[XRE]",
		"EnableProfileMigrator=0",
		"EnableExtensionManager=1",
	}
	ini := ParseIni(iniData)
	PrintIni(ini)
}

func UniqueInts(slice []int) []int {
	seen := map[int]bool{}
	unique := []int{}
	for _, x := range slice {
		if _, found := seen[x]; !found {
			unique = append(unique, x)
			seen[x] = true
		}
	}
	return unique
}

func Flatten(matrix [][]int) []int {
	slice := make([]int, 0, len(matrix)+len(matrix[0]))
	for _, innerSlice := range matrix {
		for _, x := range innerSlice {
			slice = append(slice, x)
		}
	}
	return slice
}

func Make2D(slice []int, columns int) [][]int {
	matrix := make([][]int, neededRows(slice, columns))
	for i, x := range slice {
		row := i / columns
		column := i % columns
		if matrix[row] == nil {
			matrix[row] = make([]int, columns)
		}
		matrix[row][column] = x
	}
	return matrix
}

func neededRows(slice []int, columns int) int {
	rows := len(slice) / columns
	if (len(slice))%columns != 0 {
		rows++
	}
	return rows
}

func ParseIni(lines []string) map[string]map[string]string {
	const separator = "="
	ini := make(map[string]map[string]string)
	group := "General"
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = line[1 : len(line)-1]
		} else if i := strings.Index(line, separator); i > -1 {
			key := line[:i]
			value := line[i+len(separator):]
			if _, found := ini[group]; !found {
				ini[group] = make(map[string]string)
			}
			ini[group][key] = value
		} else {
			log.Print("failed to parse line:", line)
		}
	}
	return ini
}

func PrintIni(ini map[string]map[string]string) {
	groups := make([]string, 0, len(ini))
	for group := range ini {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	for i, group := range groups {
		fmt.Printf("[%s]\n", group)
		keys := make([]string, 0, len(ini[group]))
		for key := range ini[group] {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf("%s=%s\n", key, ini[group][key])
		}
		if i+1 < len(group) {
			fmt.Println()
		}
	}
}
