package dbrepo

import (
	"context"
	"time"

	"github.com/jagottsicher/myGoWebApplication/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation stores a reservation in the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `
		insert into reservations 
			(full_name, email, phone, start_date, end_date, bungalow_id, created_at, updated_at)
		values
			($1, $2, $3, $4, $5, $6, $7, $8) returning id
	`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FullName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.BungalowID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertBungalowRestriction places a restriction in the database
func (m *postgresDBRepo) InsertBungalowRestriction(r models.BungalowRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into bungalow_restrictions
			(start_date, end_date, bungalow_id, reservation_id, created_at, updated_at, restriction_id)
		values
			($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.BungalowID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByBungalowID returns true if there is availablity for a bungalowID for a date range, false if not
func (m *postgresDBRepo) SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `
		select 
			count(id)
		from
			bungalow_restrictions
		where
			bungalow_id = $1
			$2 < end_date and $3 > start_date;
	`

	row := m.DB.QueryRowContext(ctx, query, bungalowID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}
