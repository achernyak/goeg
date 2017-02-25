package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var workers = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <image files>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	files := os.Args[1:]
	jobs := make(chan string, workers*16)
	results := make(chan string)
	done := make(chan struct{}, workers)

	go addJobs(files, jobs)
	for i := 0; i < workers; i++ {
		go doJobs(done, results, jobs)
	}
	waitAndProcessResults(done, results)
}

func addJobs(files []string, jobs chan<- string) {
	for _, filename := range files {
		jobs <- filename
	}
	close(jobs)
}

func doJobs(done chan<- struct{}, results chan<- string,
	jobs <-chan string) {
	for job := range jobs {
		if result, ok := process(job); ok {
			results <- result
		}
	}
	done <- struct{}{}
}

func waitAndProcessResults(done <-chan struct{}, results <-chan string) {
	for working := workers; working > 0; {
		select {
		case result := <-results:
			fmt.Println(result)
		case <-done:
			working--
		}
	}
DONE:
	for {
		select {
		case result := <-results:
			fmt.Println(result)
		default:
			break DONE
		}
	}
}

func process(filename string) (string, bool) {
	if info, err := os.Stat(filename); err != nil ||
		(info.Mode()&os.ModeType == 1) {
		return "", false
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", false
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return "", false
	}

	return fmt.Sprintf(`<img src="%s" width="%d" height="%d" />`,
		filepath.Base(filename), config.Width, config.Height), true
}
