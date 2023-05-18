package forms

type errors map[string][]string

// Add will add an error message to specific form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message found for a specific form field
func (e errors) Get(field string) string {
	errorString := e[field]
	if len(errorString) == 0 {
		return ""
	}
	return errorString[0]
}
