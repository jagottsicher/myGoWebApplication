package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var allTheHandlerTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"eremite", "/eremite", "GET", []postData{}, http.StatusOK},
	{"couple", "/couple", "GET", []postData{}, http.StatusOK},
	{"family", "/family", "GET", []postData{}, http.StatusOK},
	{"reservation", "/reservation", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"reservation-overview", "/reservation-overview", "GET", []postData{}, http.StatusOK},
	{"not-existing-route", "/not-existing-dummy", "GET", []postData{}, http.StatusNotFound},
}

func TestAllTheHandlers(t *testing.T) {

	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range allTheHandlerTests {
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("%s: expected %d, got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		} else {
			// here POST request handlers tests laster
		}
	}
}
