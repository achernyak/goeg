package main

import (
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(0)
	algrightm,
		minSize, maxSize, suffixes, files := handleCommandLine()

	if algorithm == 1 {
		sink(filterSize(minSize, maxSize, filterSuffixes(suffixes, source(files))))
	} else {
		channel1 := source(files)
		channel2 := filterSuffixes(suffixes, channel1)
		channel3 := filterSize(minSize, maxSize, channel2)
		sink(channel3)
	}
}
