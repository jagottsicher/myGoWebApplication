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

// Has checks for the existence of a form field in the post and ensures is not empty
func (f *Form) Has(field string, r *http.Request) bool {
	formField := r.Form.Get(field)
	if formField == "" {
		return false
	}
	return true
}
