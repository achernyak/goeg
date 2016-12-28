package main

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"strconv"
)

const (
	pageTop = `<!DOCTYPE HTML><html><head>
							<style>.error{color:#FF0000;}</style></head>
							<title>Statistics</title>
							<body><h3>Statistics</h3>
							<p>Computes basic statistics for a given list of numbers</p>`
	form = `<form action="/" method="POST">
						<input type="text" name="a" size="1"><label for="a"><i>x</i></label> +
						<input type="text" name="b" size="1"><label for="b"><i>x</i></label> +
						<input type="text" name="c" size="1"><label for="c"> </label>
						<input type="submit" name="calculate" value="Calculate">
					</form>`
	pageBottom = `</body></html>`
	anError    = `<p class="error">%s</p>`
	solution   = "<p>%s -> %s</p"
)

func main() {
	http.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	} else {
		if numbers, message, ok := processRequest(request); ok {
			question := formatQuestion(request.Form)
			x1, x2 := solve(numbers)
			answer := formatSolutions(x1, x2)
			fmt.Fprintf(writer, solution, question, answer)
		} else if message != "" {
			fmt.Fprintf(writer, anError, message)
		}
	}
	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([3]float64, string, bool) {
	var floats [3]float64
	count := 0
	for index, key := range []string{"a", "b", "c"} {
		if slice, found := request.Form[key]; found && len(slice) > 0 {
			if slice[0] != "" {
				if x, err := strconv.ParseFloat(slice[0], 64); err != nil {
					return floats, "'" + slice[0] + "' is invalid", false
				} else {
					floats[index] = x
				}
			} else {
				request.Form[key][0] = "0"
				floats[index] = 0
			}
			count++
		}
	}
	if count != 3 {
		return floats, "", false
	}
	if EqualFloat(floats[0], 0, -1) {
		return floats, "the x^2 factor may not be 0", false
	}
	return floats, "", true
}

func formatQuestion(form map[string][]string) string {
	return fmt.Sprintf("%s<i>x</i>^2 + %s<i>x</i> + %s", form["a"][0],
		form["b"][0], form["c"][0])
}

func formatSolutions(x1, x2 complex128) string {
	if EqualComplex(x1, x2) {
		return fmt.Sprintf("<i>x</x>=%f", x1)
	}
	return fmt.Sprintf("<i>x</i>=%f or <i>x</i>=%f", x1, x2)
}

func solve(floats [3]float64) (complex128, complex128) {
	a, b, c := complex(floats[0], 0), complex(floats[1], 0),
		complex(floats[2], 0)
	root := cmplx.Sqrt(cmplx.Pow(b, 2) - (4 * a * c))
	x1 := (-b + root) / (2 * a)
	x2 := (-b - root) / (2 * a)
	return x1, x2
}

func EqualFloat(x, y, limit float64) bool {
	if limit <= 0.0 {
		limit = math.SmallestNonzeroFloat64
	}
	return math.Abs(x-y) <= (limit * math.Min(math.Abs(x), math.Abs(y)))
}

func EqualComplex(x, y complex128) bool {
	return EqualFloat(real(x), real(y), -1) &&
		EqualFloat(imag(x), imag(y), -1)
}
