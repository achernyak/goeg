package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	externalLinkRx *regexp.Regexp
	addChannel     chan string
	queryChannel   chan string
	seenChannel    chan bool
)

func init() {
	externalLinkRx = regexp.MustCompile("^(http|ftp|mailto):")
	addChannel = make(chan string)
	queryChannel = make(chan string)
	seenChannel = make(chan bool)
}

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		log.Fatalf("usage: %s url\n", filepath.Base(os.Args[0]))
	}
	href := os.Args[1]
	if !strings.HasPrefix(href, "http://") {
		href = "http://" + href
	}
	url, err := url.Parse(href)
	if err != nil {
		log.Fatalln("- failed to read url:", err)
	}
	prepareMap()
	checkPage(href, "http://"+url.Host)
}
