package dbrepo

import (
	"errors"
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
	return false, nil
}

// SearchAvailabilityByDatesForAllBungalows returns a slice of available bungalows, if any for a queried date range
func (m *testDBRepo) SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error) {
	var bungalows []models.Bungalow

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
