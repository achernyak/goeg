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
	decimals = 3
	pageTop  = `<!DOCTYPE HTML><html><head>
							<style>.error{color:#FF0000;}</style></head>
							<title>Quadratic Equation Solver</title>
							<body><h3>Quadratic Equation Solver</h3>
							<p>solves equations of the form
								a<i>x</i>^2 + b<i>x</i> + c</p>`
	form = `<form action="/" method="POST">
						<input type="text" name="a" size="1"><label for="a"><i>x</i></label> +
						<input type="text" name="b" size="1"><label for="b"><i>x</i></label> +
						<input type="text" name="c" size="1"><label for="c"> </label>
						<input type="submit" name="calculate" value="Calculate">
					</form>`
	pageBottom   = `</body></html>`
	anError      = `<p class="error">%s</p>`
	solution     = "<p>%s -> %s</p"
	oneSolution  = "<i>x</i>=%s"
	twoSolutions = "<i>x</i>=%s or <i>x</i>=%s"
	noSolution   = "<i>there are no solutions</i>"
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
	result := formatSignAndNumber("", form["a"][0], "<i>x</i>^2")
	result += formatSignAndNumber(" ", form["b"][0], "<i>x</i>")
	result += formatSignAndNumber(" ", form["c"][0], "")
	return result
}

func formatSignAndNumber(signPad, number, suffix string) string {
	if number == "" || number == "0" || number == "0.0" {
		return ""
	}
	var sign string
	if signPad != "" {
		sign = signPad + "+" + signPad
	}
	if number[0] == '-' {
		sign = signPad + "-" + signPad
		number = number[1:]
	}
	if suffix != "" && number == "1" {
		return sign + suffix
	}
	return sign + number + suffix
}

func formatSolutions(x1, x2 complex128) string {
	exactlyOneSolution := false
	if cmplx.IsNaN(x1) && cmplx.IsNaN(x2) {
		return noSolution
	}
	if cmplx.IsNaN(x1) {
		exactlyOneSolution = true
		x1 = x2
	} else if cmplx.IsNaN(x2) || EqualComplex(x1, x2) {
		exactlyOneSolution = true
	}
	if exactlyOneSolution {
		return fmt.Sprintf(oneSolution, formatComplex(x1))
	}
	return fmt.Sprintf(twoSolutions, formatComplex(x1), formatComplex(x2))
}

func formatComplex(x complex128) string {
	if EqualFloat(imag(x), 0, -1) {
		return fmt.Sprintf("%.*f", decimals, real(x))
	}
	return fmt.Sprintf("%.*f", decimals, x)
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
