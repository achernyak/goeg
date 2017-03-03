package linkutil

import (
	"bufio"
	"io"
	"linkcheck/linkutil"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestLinksFromReader(t *testing.T) {
	file, err := os.Open("index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	links, err := linkutil.LinksFromReader(file)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(links)
	file, err = os.Open("index.links")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			lines = append(lines, line[:len(line)-1])
		}
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
	}
	sort.Strings(lines)

	if !reflect.DeepEqual(links, lines) {
		for i := 0; i < len(links); i++ {
			if i < len(lines) {
				if links[i] != lines[i] {
					t.Fatalf("%q != %q", links[i], lines[i])
				}
			} else {
				t.Fatalf("found more links than lines, e.g.: %q", links[i])
			}
		}
		t.Fatalf("found fewer links than lines (%d vs. %d)", len(links),
			len(lines))
	}
}
