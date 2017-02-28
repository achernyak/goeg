package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		log.Fatalf("usage: %s utf-16-in.txt [>]utf-8-out.txt\n",
			filepath.Base(os.Args[0]))
		return
	}
	var err error
	var infile *os.File
	if infile, err = os.Open(os.Args[1]); err != nil {
		log.Fatalln(err)
	}
	defer infile.Close()
	outfile := os.Stdout
	if len(os.Args) > 2 {
		if outfile, err = os.Create(os.Args[2]); err != nil {
			log.Fatalln(err)
		}
		defer outfile.Close()
	}
	if err := utf16toutf8(infile, outfile); err != nil {
		log.Fatalln(err)
	}
}
