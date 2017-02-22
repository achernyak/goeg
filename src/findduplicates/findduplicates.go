package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type fileInfo struct {
	sha1 []byte
	size int64
	path string
}

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

func findDuplicates(infoChan chan fileInfo, dirname string) {
	waiter := &sync.WaitGroup{}
	filepath.Walk(dirname, makeWalkFunc(infoChan, waiter))
	waiter.Wait()
	close(infoChan)
}

const maxSizeOfSmallFile = 1024 * 32

func makeWalkFunc(infoChan chan fileInfo,
	waiter *sync.WaitGroup) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Size() > 0 &&
			(info.Mode()&os.ModeType == 0) {
			if info.Size() < maxSizeOfSmallFile ||
				runtime.NumGoroutine() > maxGoroutines {
				processFile(path, info, infoChan, nil)
			} else {
				waiter.Add(1)
				go processFile(path, info, infoChan,
					func() { waiter.Done() })
			}
		}
		return nil
	}
}

func processFile(filename string, info os.FileInfo,
	infoChan chan fileInfo, done func()) {
	if done != nil {
		defer done()
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Println("error:", err)
		return
	}
	defer file.Close()
	hash := sha1.New()
	if size, err := io.Copy(hash, file); size != info.Size() || err != nil {
		if err != nil {
			log.Println("error:", err)
		} else {
			log.Println("error: failed to read the whole file:", filename)
		}
		return
	}
	infoChan <- fileInfo{hash.Sum(nil), info.Size(), filename}
}
