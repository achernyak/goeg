package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const maxGoroutines = 100

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <path>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	infoChan := make(chan fileInfo, maxGoroutines*2)
	go findDuplicates(infoChan, os.Args[1])
	pathData := mergeResults(infoChan)
	outputResults(pathData)
}
