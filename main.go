package main

import (
	"errors"
	"fmt"
	"net/http"
)

const portNumber = ":8080"

// Home is the handler for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
}

// About is the handler for the about page
func About(w http.ResponseWriter, r *http.Request) {
	sum := addValues(2, 2)
	_, _ = fmt.Fprintf(w, fmt.Sprintf("This is the about page and 2 + 2 is %d", sum))
}

// Divide divides one value into another and returns message with result
func Divide(w http.ResponseWriter, r *http.Request) {
	x := 100.0
	y := 0.0
	f, err := divideValues(x, y)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: Division by zero is not a valid operation. Error returned: %s", err)
		return
	}

	_, _ = fmt.Fprintf(w, fmt.Sprintf("%f divided by %f is %f", x, y, f))
}

// divideValues divides two floats x and y, and returns the quotient and a value of type error
func divideValues(x, y float64) (float64, error) {
	if y == 0 {
		err := errors.New("Divisor is zero!")
		return 0, err
	}
	result := x / y
	return result, nil
}

// addValues adds two ints x and y, and returns the sum
func addValues(x, y int) int {
	return x + y
}

// main is the main function
func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/divide", Divide)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
