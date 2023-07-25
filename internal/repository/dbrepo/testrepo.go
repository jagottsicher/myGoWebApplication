package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/jagottsicher/myGoWebApplication/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation stores a reservation in the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.BungalowID == 99 {
		return 0, errors.New("some error")
	}

	return 1, nil
}

// InsertBungalowRestriction places a restriction in the database
func (m *testDBRepo) InsertBungalowRestriction(r models.BungalowRestriction) error {
	if r.BungalowID == 999 {
		return errors.New("just because")
	}

	return nil
}

// SearchAvailabilityByDatesByBungalowID returns true if there is availablity for a bungalowID for a date range, false if not
func (m *testDBRepo) SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowID int) (bool, error) {
	// set up a test time
	layout := "2006-01-02"
	str := "2036-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2038-01-01 as start
	testDateToFail, err := time.Parse(layout, "2038-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2036-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

// SearchAvailabilityByDatesForAllBungalows returns a slice of available bungalows, if any for a queried date range
func (m *testDBRepo) SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error) {
	var bungalows []models.Bungalow

	// if the start date is after 2036-12-31, then return empty slice,
	// indicating no bungalows are available;
	layout := "2006-01-02"
	str := "2036-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2038-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return bungalows, errors.New("some error")
	}

	if start.After(t) {
		return bungalows, nil
	}

	// otherwise, put an entry into the slice, indicating that some bungalow is
	// available for search dates
	bungalow := models.Bungalow{
		ID: 1,
	}
	bungalows = append(bungalows, bungalow)

	return bungalows, nil
}

// GetBungalowByID gets a bungalow by id
func (m *testDBRepo) GetBungalowByID(id int) (models.Bungalow, error) {
	var bungalow models.Bungalow
	if id > 3 {
		return bungalow, errors.New("an error occured")
	}

	return bungalow, nil
}
