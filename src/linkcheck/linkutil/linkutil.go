package linkutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

var hrefRx *regexp.Regexp

func init() {
	hrefRx = regexp.MustCompile(`<a[^>]+href=['"]?([^'">]+)['"]?`)
}

func LinksFromURL(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %s", err)
	}
	links, err := LinksFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse page: %s", err)
	}
	return links, nil
}

func LinksFromReader(reader io.Reader) ([]string, error) {
	html, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	uniqueLinks := make(map[string]bool)
	for _, submatch := range hrefRx.FindAllSubmatch(html, -1) {
		uniqueLinks[string(submatch[1])] = true
	}
	links := make([]string, len(uniqueLinks))
	i := 0
	for link := range uniqueLinks {
		links[i] = link
		i++
	}
	return links, nil
}
