package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"post-reservation", "/reservation", "POST", []postData{
		{key: "startingEnd", value: "2023-02-23"},
		{key: "endingEnd", value: "2023-02-25"},
	}, http.StatusOK},
	{"post-reservation-json", "/reservation-json", "POST", []postData{
		{key: "startingEnd", value: "2023-02-23"},
		{key: "endingEnd", value: "2023-02-25"},
	}, http.StatusOK},
	{"post-make-reservation", "/make-reservation", "POST", []postData{
		{key: "full_name", value: "Ricky Spanish"},
		{key: "email", value: "me@writing-an-email.to"},
		{key: "phone", value: "555-123-4567"},
	}, http.StatusOK},
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
			values := url.Values{}

			for _, param := range test.params {
				values.Add(param.key, param.value)
			}

			response, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("%s: expected %d,  got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
