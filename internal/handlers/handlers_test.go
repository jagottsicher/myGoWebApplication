package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jagottsicher/myGoWebApplication/internal/driver"
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
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"eremite", "/eremite", "GET", http.StatusOK},
	{"couple", "/couple", "GET", http.StatusOK},
	{"family", "/family", "GET", http.StatusOK},
	{"reservation", "/reservation", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"not-existing-route", "/not-existing-dummy", "GET", http.StatusNotFound},
}

// TestNewRepo test NewRepo only by a package reflect
func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

// TestAllTheHandlers tests the trivial get-request handlers only
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
		}
	}
}

// TestRepository_PostReservation tests the PostReservation post-request handler
func TestRepository_PostReservation(t *testing.T) {

	// case #1: bungalows are not available

	// create request body
	postedData := url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "2037-01-02")

	// create  request
	req, _ := http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr := httptest.NewRecorder()

	// create handler to function to test
	handler := http.HandlerFunc(Repo.PostReservation)

	// make request to handler
	handler.ServeHTTP(rr, req)

	// since we have no bungalows available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability when no bungalows available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case #2: bungalows are available

	// this time, we specify a start date before 2040-01-01, which will give us
	// a non-empty slice, indicating that bungalows are available
	postedData = url.Values{}
	postedData.Add("start", "2036-01-01")
	postedData.Add("end", "2036-01-02")

	// create request
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.PostReservation)

	// make request to handler
	handler.ServeHTTP(rr, req)

	// since we have bungalows available, we expect to get status http.StatusOK
	if rr.Code != http.StatusOK {
		t.Errorf("Post availability when bungalows are available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// case #3: empty post body

	// create request with a nil body, so parsing form fails
	req, _ = http.NewRequest("POST", "/reservation", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.PostReservation)

	// make request to handler
	handler.ServeHTTP(rr, req)

	// bungalows available, expected status code http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #4: start date in wrong format

	// start date in wrong format
	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2037-01-02")

	// create request
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.PostReservation)

	// make request to handler
	handler.ServeHTTP(rr, req)

	// bungalows available, expected status code http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #5: end date in wrong format

	// end date in wrong format
	postedData = url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "invalid")

	// create request
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.PostReservation)

	// make request to handler
	handler.ServeHTTP(rr, req)

	// bungalows available, expected status code http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #6: database query fails

	// start date of 2038-01-01, expects testdb repo to return an error
	postedData = url.Values{}
	postedData.Add("start", "2038-01-01")
	postedData.Add("end", "2038-01-02")

	// create request
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.PostReservation)

	// make request to handler
	handler.ServeHTTP(rr, req)

	// bungalows available, expected status code http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_MakeReservation tests the MakeReservation get-request handler
func TestRepository_MakeReservation(t *testing.T) {

	// case #1: all okay

	// creating reservation data for test purpose

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

	// case #2: without a reservation in session (reset everything)

	req, _ = http.NewRequest("GET", "/make-reservation", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #3: test error returned from database query function

	req, _ = http.NewRequest("GET", "/make-reservation", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.BungalowID = 99
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// TestRepository_PostMakeReservation tests the MakeReservation post-request handler
func TestRepository_PostMakeReservation(t *testing.T) {

	// case #1: reservation works fine

	postedData := url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2037-01-01")
	ed, _ := time.Parse(layout, "2037-01-02")
	bungalowId, _ := strconv.Atoi("1")

	reservation := models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case #2: missing post body

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #3: missing session data

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code for missing session data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #4: invalid/insufficient data

	postedData = url.Values{}
	postedData.Add("full_name", "P")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("PostMakeReservation handler returned wrong response code invalid/insufficient data: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// case #5:  failure inserting reservation into database

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("99")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #6: failure to inserting restriction into database

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("999")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_ReservationJSON tests the ReservationJSON POST-request handler
func TestRepository_ReservationJSON(t *testing.T) {

	// case #1: bungalow ID invalid

	postedData := url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "2037-01-02")
	postedData.Add("bungalow_id", "invalid")

	// create request
	req, _ := http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationJSON handler unexpectedly seems able parsing invalid bungalow ID: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #2: start date invalid

	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2037-01-02")
	postedData.Add("bungalow_id", "1")

	// create request
	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationJSON handler unexpectedly seems able parsing invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #3: end date invalid

	postedData = url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "invalid")
	postedData.Add("bungalow_id", "1")

	// create request
	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationJSON handler unexpectedly seems able parsing invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #4: bungalow not available

	postedData = url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "2037-01-02")
	postedData.Add("bungalow_id", "1")

	// create request
	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	// no bungalows available, expected status code http.StatusSeeOther
	//  parse JSON and receive response
	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("can't parse json")
	}

	// start date > 2036-12-31, expecting no availability
	if j.OK {
		t.Error("got availability, unexpected none")
	}

	// case #5: bungalow available

	// create request body
	postedData = url.Values{}
	postedData.Add("start", "2036-01-01")
	postedData.Add("end", "2036-01-02")
	postedData.Add("bungalow_id", "1")

	// create  request
	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.ReservationJSON)

	// make request to handler
	handler.ServeHTTP(rr, req)

	// parse JSON and receive response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// start date < 2036-12-31, expecting availability
	if !j.OK {
		t.Error("got no availability but expected in handler ReservationJSON")
	}

	// case #6:  no request body

	// create  request
	req, _ = http.NewRequest("POST", "/reservation-json", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.ReservationJSON)

	// make request to handler
	handler.ServeHTTP(rr, req)

	//  parse JSON and receive response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// start date < 2036-12-31, expecting availability
	if j.OK || j.Message != "Internal server error" {
		t.Error("got availability with empty request body")
	}

	// case #7: database error

	// create request body
	postedData = url.Values{}
	postedData.Add("start", "2038-01-01")
	postedData.Add("end", "2038-01-02")
	postedData.Add("bungalow_id", "1")

	// create request
	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create response recorder
	rr = httptest.NewRecorder()

	// create handler to function to test
	handler = http.HandlerFunc(Repo.ReservationJSON)

	// make request to handler
	handler.ServeHTTP(rr, req)

	//  parse JSON and receive response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// start date = 2038-01-01, expecting no availability
	if j.OK || j.Message != "Error querying database" {
		t.Error("got availability, unexpected because database returned an error")
	}
}

// TestRepository_ReservationOverview tests the ReservationOverview request handler
func TestRepository_ReservationOverview(t *testing.T) {

	// case #1: reservation in session

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-overview", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationOverview)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationOverview handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// case #2: bungalow not in database

	reservation = models.Reservation{
		BungalowID: 99,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ = http.NewRequest("GET", "/reservation-overview", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.ReservationOverview)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationOverview handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #3: reservation not in session

	req, _ = http.NewRequest("GET", "/reservation-overview", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationOverview)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationOverview handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

// TestRepository_ChooseBungalow tests the ChooseBungalow handler
func TestRepository_ChooseBungalow(t *testing.T) {

	// case #1: reservation in session

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-bungalow/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the RequestURI on the request to get ID from the URL
	req.RequestURI = "/choose-bungalow/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case #2: reservation not in session

	req, _ = http.NewRequest("GET", "/choose-bungalow/1", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-bungalow/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #3: missing url parameter, or malformed parameter

	req, _ = http.NewRequest("GET", "/choose-bungalow/fish", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.RequestURI = "/choose-bungalow/fish"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_BookBungalow tests the BookBungalow handler
func TestRepository_BookBungalow(t *testing.T) {

	// case #1: database works

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/book-bungalow?s=2038-01-01&e=2038-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case #2: database failed

	req, _ = http.NewRequest("GET", "/book-bungalow?s=2036-01-01&e=2036-01-02&id=99", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}

// loginTests is the data for the Login handler tests
var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"patrick@bikini-bottom.ocean",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminPostShowReservationTests = []struct {
	name                 string
	url                  string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name: "valid-data-from-new",
		url:  "/admin/reservations/new/1/show",
		postedData: url.Values{
			"full_name": {"Stan Smith"},
			"email":     {"stan-the-man@cia.com"},
			"phone":     {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-all",
		url:  "/admin/reservations/all/1/show",
		postedData: url.Values{
			"full_name": {"Stan Smith"},
			"email":     {"stan-the-man@cia.com"},
			"phone":     {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-all",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-cal",
		url:  "/admin/reservations/calendar/1/show",
		postedData: url.Values{
			"full_name": {"Stan Smith"},
			"email":     {"stan-the-man@cia.com"},
			"phone":     {"555-555-5555"},
			"year":      {"2024"},
			"month":     {"02"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-calendar?y=2024&m=02",
		expectedHTML:         "",
	},
}

// TestAdminPostShowReservation tests the AdminPostReservation handler
func TestAdminPostShowReservation(t *testing.T) {
	for _, e := range adminPostShowReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/user/login", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/user/login", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostShowReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminPostReservationCalendarTests = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
	blocks               int
	reservations         int
}{
	{
		name: "cal",
		postedData: url.Values{
			"year":  {time.Now().Format("2006")},
			"month": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2006-01-2")): {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "cal-blocks",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		blocks:               1,
	},
	{
		name:                 "cal-res",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		reservations:         1,
	},
}

func TestPostReservationCalendar(t *testing.T) {
	for _, e := range adminPostReservationCalendarTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		now := time.Now()
		bm := make(map[string]int)
		rm := make(map[string]int)

		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2006-01-2")] = 0
			bm[d.Format("2006-01-2")] = 0
		}

		if e.blocks > 0 {
			bm[firstOfMonth.Format("2006-01-2")] = e.blocks
		}

		if e.reservations > 0 {
			rm[lastOfMonth.Format("2006-01-2")] = e.reservations
		}

		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostReservationsCalendar)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

	}
}

var adminProcessReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "process-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "process-reservation-back-to-cal",
		queryParams:          "?y=2024&m=02",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminProcessReservation(t *testing.T) {
	for _, e := range adminProcessReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminDeleteReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "delete-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "delete-reservation-back-to-cal",
		queryParams:          "?y=2024&m=02",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, e := range adminDeleteReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}
