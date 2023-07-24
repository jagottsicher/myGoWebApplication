package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

	// test case without a reservation in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test error returned from database query function
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
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

// TestRepository_PostMakeReservation tests the MakeReservation get-request handler
func TestRepository_PostMakeReservation(t *testing.T) {

	reqBody := "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2037-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=peter@griffin.family")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "bungalow_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2037-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=peter@griffin.family")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "bungalow_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid end date
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=peter@griffin.family")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "bungalow_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid bungalow id
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2037-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=peter@griffin.family")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "bungalow_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code for invalid bungalow id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid/insufficient data
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2037-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "full_name=P")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=peter@griffin.family")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "bungalow_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code for invalid bungalow id: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// failure inserting reservation into database
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2037-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=peter@griffin.family")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "bungalow_id=99")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// failure to inserting restriction into database
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2037-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=peter@griffin.family")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "bungalow_id=999")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
