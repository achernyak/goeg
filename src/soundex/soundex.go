package main

import (
	"log"
	"net/http"
)

const (
	pageTop = `<!DOCTYPE HTML><html><head>
			<style>.error{color:#FF0000;}</style></head>
			<title>Soundex</title><body><h3>Soundex</h3>
			<p>Compute soundex codes for a list of names.</p>`
	form = `<form action="/" method="POST">
					<label for="names">Names (comma or space-separated):</label><br />
					<input type="text" name="names" size="30"><br />
					<input type="submit" name="compute" value="Compute">
					</form>`
	pageBottom = `</body></html>`
	error      = `<p class="error">%s</p>`
)

var digitForLetter = []rune{
	0, 1, 2, 3, 0, 1, 2, 0, 0, 2, 2, 4, 5,
	// A  B  C  D  E  F  G  H  I  J  K  L  M
	5, 0, 1, 2, 6, 2, 3, 0, 1, 0, 2, 0, 2}

// N  O  P  Q  R  S  T  U  V  W  X  Y  Z

var testCases map[string]string

func main() {
	http.HandleFunc("/", homePage)
	var ok bool
	if testCases, ok = populateTestCases("soundex-test-data.txt"); ok {
		http.HandleFunc("/test", testPage)
	}
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}
