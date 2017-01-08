package main

import "fmt"

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
