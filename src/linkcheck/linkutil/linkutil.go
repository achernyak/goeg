package linkutil

import "regexp"

var hrefRx *regexp.Regexp

func init() {
	hrefRx = regexp.MustCompile(`<a[^>]+href=['"]?([^'">]+)['"]?`)
}
