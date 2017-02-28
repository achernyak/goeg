package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"unicode/utf16"
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

func utf16toutf8(infile, outfile *os.File) error {
	writer := bufio.NewWriter(outfile)
	defer writer.Flush()
	bom := make([]byte, 2)
	if _, err := infile.Read(bom); err != nil {
		return err
	}
	var byteOrder binary.ByteOrder = binary.LittleEndian
	if bom[0] == 0xFE && bom[1] == 0xFF {
		byteOrder = binary.BigEndian
	} else if bom[0] != 0xFF && bom[1] != 0xFE {
		return errors.New("missing byte order mark")
	}
	for {
		var c uint16
		if err := binary.Read(infile, byteOrder, &c); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if _, err := writer.WriteString(
			string(utf16.Decode([]uint16{c}))); err != nil {
			return err
		}
	}
}
