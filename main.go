package main

import (
	"fmt"
	"net/http"
)

// Home is the handler for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
}

// About is the handler for the about page
func About(w http.ResponseWriter, r *http.Request) {
	owner, saying := getData()
	sum := addValues(2, 2)
	fmt.Fprintf(w, fmt.Sprintf("This is the about page of %s, \nI like to say %s, \nand as a side note, 2 + 2 is %d", owner, saying, sum))
}

// addValues adds two ints x and y, and returns the sum
func addValues(x, y int) int {
	return x + y
}

// getData returns a name and a saying
func getData() (string, string) {
	o := "Rick Sanchez"
	s := "Wubba Lubba Dup Dup!"
	return o, s
}

// main is the main function
func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	_ = http.ListenAndServe(":8080", nil)
}
