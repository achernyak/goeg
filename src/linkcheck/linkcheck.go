package main

import (
	"fmt"
	"linkcheck/linkutil"
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

func prepareMap() {
	go func() {
		seen := make(map[string]bool)
		for {
			select {
			case url := <-addChannel:
				seen[url] = true
			case url := <-queryChannel:
				_, found := seen[url]
				seenChannel <- found
			}
		}
	}()
}

func alreadySeen(url string) bool {
	queryChannel <- url
	if <-seenChannel {
		return true
	}
	addChannel <- url
	return false
}

func checkPage(url, site string) {
	if alreadySeen(url) {
		return
	}
	links, err := linkutil.LinksFromURL(url)
	if err != nil {
		log.Println("-", err)
		return
	}
	fmt.Println("+ read", url)
	done := make(chan bool, len(links))
	defer close(done)
	pending := 0
	var messages []string
	for _, link := range links {
		pending += processLink(link, site, url, &messages, done)
	}
	if len(messages) > 0 {
		fmt.Println("+ links on", url)
		for _, message := range messages {
			fmt.Println("  ", message)
		}
	}
	for i := 0; i < pending; i++ {
		<-done
	}
}

func processLink(link, site, url string, messages *[]string,
	done chan<- bool) int {
	localAndParsable, link := classifyLink(link, site)
	if localAndParsable {
		go func() {
			checkPage(link, site)
			done <- true
		}()
		return 1
	}
	if message := checkExists(link, url); message != "" {
		*messages = append(*messages, message)
	}
	return 0
}
