package forms

import (
	"net/http"
	"net/url"
)

// Form is a type holding a genaral form struct including an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New is a function to initialize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid returns false in case of errors, otherwise true
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Has checks for the existence of a form field in the post and ensures is not empty
func (f *Form) Has(field string, r *http.Request) bool {
	formField := r.Form.Get(field)
	if formField == "" {
		f.Errors.Add(field, "This field cannot be empty.")
		return false
	}
	return true
}
