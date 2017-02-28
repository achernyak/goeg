package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func unpackZip(filename string) (err error) {
	var reader *zip.ReadCloser
	if reader, err = zip.OpenReader(filename); err != nil {
		return err
	}
	defer reader.Close()
	for _, zipFile := range reader.Reader.File {
		filename := sanitizedName(zipFile.Name)
		if strings.HasSuffix(zipFile.Name, "/") ||
			strings.HasSuffix(zipFile.Name, "\\") {
			if err = os.MkdirAll(filename, 0755); err != nil {
				return err
			}

		} else {
			if err = unpackZippedFile(filename, zipFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func unpackZippedFile(filename string, zipFile *zip.File) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()
	reader, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()
	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	if filename == zipFile.Name {
		fmt.Println(filename)
	} else {
		fmt.Printf("%s [%s]\n", filename, zipFile.Name)
	}
	return nil
}

func unpackTarFiles(reader *tar.Reader) error {
	for {
		header, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		filename := sanitizedName(header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filename, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err = unpackTarFile(filename, header.Name, reader); err != nil {
				return err
			}
		}
	}
}

func unpackTarFile(filename, tarFilename string, reader *tar.Reader) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()
	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	if filename == tarFilename {
		fmt.Println(filename)
	} else {
		fmt.Printf("%s [%s]\n", filename, tarFilename)
	}
	return nil
}

func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	filename = filepath.ToSlash(filename)
	filename = strings.TrimLeft(filename, "/.")
	return strings.Replace(filename, "../", "", -1)
}
