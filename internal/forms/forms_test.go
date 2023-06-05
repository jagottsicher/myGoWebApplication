package forms

import (
	"testing"
)

func TestForm_Valid(t *testing.T) {

	if !isValid {
		t.Error("expected valid, got invalid")
	}
}

func TestForm_Required(t *testing.T) {

	if form.Valid() {
		t.Error("required fields missing, expected invalid, got valid")
	}

	if !form.Valid() {
		t.Error("reports no required fields, but does")
	}
}

func TestForm_Has(t *testing.T) {

	if has {
		t.Error("reports an existing field, but does not")
	}

	if !has {
		t.Error("reports not existing field, when it should")
	}
}

func TestForm_MinLength(t *testing.T) {

	if form.Valid() {
		t.Error("form shows mininum length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should report an error but does not")
	}

	if form.Valid() {
		t.Error("shows minimum length of 100 met but is shorter")
	}

	if !form.Valid() {
		t.Error("shows minimum length of 1 is not met when it is")
	}

	if isError != "" {
		t.Error("should not have error but reports one")
	}

}

func TestForm_IsEmail(t *testing.T) {

	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	if form.Valid() {
		t.Error("got valid for invalid email address")
	}
}
