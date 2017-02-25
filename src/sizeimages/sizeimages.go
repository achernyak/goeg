package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

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

	files := os.Args[1:]
	jobs := make(chan string, workers*16)
	done := make(chan struct{}, workers)
	go addJobs(files, jobs)
	for i := 0; i < workers; i++ {
		go doJobs(done, jobs)
	}
	waitUntil(done)
}

func addJobs(files []string, jobs chan<- string) {
	for _, filename := range files {
		suffix := strings.ToLower(filepath.Ext(filename))
		if suffix == ".html" || suffix == ".htm" {
			jobs <- filename
		}
	}
	close(jobs)
}

func doJobs(done chan<- struct{}, jobs <-chan string) {
	for job := range jobs {
		sizeImages(job)
	}
	done <- struct{}{}
}

func waitUntil(done <-chan struct{}) {
	for i := 0; i < workers; i++ {
		<-done
	}
}
