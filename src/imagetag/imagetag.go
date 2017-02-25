package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var workers = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <image files>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	files := commandLineFiles(os.Args[1:])
	jobs := make(chan string, workers*16)
	results := make(chan string)
	done := make(chan struct{}, workers)

	go addJobs(files, jobs)
	for i := 0; i < workers; i++ {
		go doJobs(done, results, jobs)
	}
	waitAndProcessResults(done, results)
}
