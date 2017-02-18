package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

var workers = runtime.NumCPU()

type Result struct {
	filename string
	lino     int
	line     string
}

type Job struct {
	filename string
	results  chan<- Result
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) < 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <regexp> <files>\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	if lineRx, err := regexp.Compile(os.Args[1]); err != nil {
		log.Fatalf("invalid regexp: %s\n", err)
	} else {
		grep(lineRx, os.Args[2:])
	}
}

func grep(lineRx *regexp.Regexp, filenames []string) {
	jobs := make(chan Job, workers)
	results := make(chan Result, minimum(1000, len(filenames)))
	done := make(chan struct{}, workers)

	go addJobs(jobs, filenames, results)
	for i := 0; i < workers; i++ {
		go doJobs(done, lineRx, jobs)
	}
	go awaitCompletion(done, results)
	processResults(results)
}

func addJobs(jobs chan<- Job, filenames []string, results chan<- Result) {
	for _, filename := range filenames {
		jobs <- Job{filename, results}
	}
	close(jobs)
}

func doJobs(done chan<- struct{}, lineRx *regexp.Regexp, jobs <-chan Job) {
	for job := range jobs {
		job.Do(lineRx)
	}
	done <- struct{}{}
}

func awaitCompletion(done <-chan struct{}, results chan Result) {
	for i := 0; i < workers; i++ {
		<-done
	}
	close(results)
}

func processResults(results <-chan Result) {
	for result := range results {
		fmt.Printf("%s:%d:%s\n", result.filename, result.lino, result.line)
	}
}

func minimum(x int, ys ...int) int {
	for _, y := range ys {
		if y < x {
			x = y
		}
	}
	return x
}
