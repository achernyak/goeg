package main

import "runtime"

var workers = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	lines := make(chan string, workers*4)
	done := make(chan struct{}, workers)
	pageMap := safemap.New()
	go readLines(os.Args[1], lines)
	processLines(done, pageMap, lines)
	waitUntil(done)
	showResults(pageMap)
}
