package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jagottsicher/myGoWebApplication/internal/models"
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
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"eremite", "/eremite", "GET", []postData{}, http.StatusOK},
	// {"couple", "/couple", "GET", []postData{}, http.StatusOK},
	// {"family", "/family", "GET", []postData{}, http.StatusOK},
	// {"reservation", "/reservation", "GET", []postData{}, http.StatusOK},
	// {"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// {"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// // {"reservation-overview", "/reservation-overview", "GET", []postData{}, http.StatusOK},
	// {"not-existing-route", "/not-existing-dummy", "GET", []postData{}, http.StatusNotFound},
	// {"post-reservation", "/reservation", "POST", []postData{
	// 	{key: "startingEnd", value: "2023-02-23"},
	// 	{key: "endingEnd", value: "2023-02-25"},
	// }, http.StatusOK},
	// {"post-reservation-json", "/reservation-json", "POST", []postData{
	// 	{key: "startingEnd", value: "2023-02-23"},
	// 	{key: "endingEnd", value: "2023-02-25"},
	// }, http.StatusOK},
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "full_name", value: "Ricky Spanish"},
	// 	{key: "email", value: "me@writing-an-email.to"},
	// 	{key: "phone", value: "555-123-4567"},
	// }, http.StatusOK},
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

// TestRepository_MakeReservation tests the MakeReservation get-request handler
func TestRepository_MakeReservation(t *testing.T) {

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// rr means "request recorder" and is an initialised response recorder for http requests builtin the test
	// basically "faking" a client and to provide a valid request/response-cycle during tests
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	// log.Println(reservation)

	// turning handler into a function
	handler := http.HandlerFunc(Repo.MakeReservation)

	// calling handler to test as a function
	handler.ServeHTTP(rr, req)

	// the test itself as a condition to test
	if rr.Code != http.StatusOK {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
