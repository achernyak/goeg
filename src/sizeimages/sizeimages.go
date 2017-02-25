package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var workers = runtime.NumCPU()

const (
	widthAttr  = "width="
	heightAttr = "height="
)

var (
	imageRx *regexp.Regexp
	srcRx   *regexp.Regexp
)

func init() {
	imageRx = regexp.MustCompile(`<[iI][mM][gG][^>]+>`)
	srcRx = regexp.MustCompile(`src=["']([^"']+)['"]`)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <html files>\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	files := commandLineFiles(os.Args[1:])
	jobs := make(chan string, workers*16)
	done := make(chan struct{}, workers)
	go addJobs(files, jobs)
	for i := 0; i < workers; i++ {
		go doJobs(done, jobs)
	}
	waitUntil(done)
}
