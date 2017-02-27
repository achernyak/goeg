package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s archive.{.zip,tar,tar.gz,tar.bz2}\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	filename := os.Args[1]
	if !validSuffix(filename) {
		log.Fatalln("unrecognized archive suffix")
	}
	if err := unpackArchive(filename); err != nil {
		log.Fatalln(err)
	}
}

func validSuffix(filename string) bool {
	for _, suffix := range []string{".zip", ".tar", ".tar.gz", ".tar.bz2"} {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}

func unpackArchive(filename string) error {
	if strings.HasSuffix(filename, ".zip") {
		return unpackZip(filename)
	}
	return unpackTar(filename)
}
