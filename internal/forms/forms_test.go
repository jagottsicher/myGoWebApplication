package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/an-url", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("expected valid, got invalid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/an-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("required fields missing, expected invalid, got valid")
	}

	postedData := url.Values{}
	postedData.Add("a", "a value")
	postedData.Add("b", "another value")
	postedData.Add("c", "yet another value")

	r, _ = http.NewRequest("POST", "/an-url", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("reports no required fields, but does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("reports an existing field, but does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("reports not existing field, when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows mininum length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should report an error but does not")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some_value")

	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows minimum length of 100 met but is shorter")
	}

	postedData = url.Values{}
	postedData.Add("some_other_field", "easy as abc")

	form = New(postedData)

	form.MinLength("some_other_field", 1)
	if !form.Valid() {
		t.Error("shows minimum length of 1 is not met when it is")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have error but reports one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "me@here.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")
	}
}
