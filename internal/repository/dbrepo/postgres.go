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
		fromw
			bungalow_restrictions
		where
			bungalow_id = $1
			$2 <= end_date and $3 >= start_date;
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

// SearchAvailabilityByDatesForAllBungalows returns a slice of available bungalows, if any for a queried date range
func (m *postgresDBRepo) SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bungalows []models.Bungalow

	query := `
		select 
			b.id, b.bungalow_name
		from
			bungalows b 
		where b.id not in 
			(select 
				bungalow_id
			from
				bungalow_restrictions br
			where 
			$1 <= br.end_date and $2 >= br.start_date
			);
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return bungalows, err
	}

	for rows.Next() {
		var bungalow models.Bungalow
		err := rows.Scan(
			&bungalow.ID,
			&bungalow.BungalowName,
		)
		if err != nil {
			return bungalows, err
		}

		bungalows = append(bungalows, bungalow)
	}

	if err = rows.Err(); err != nil {
		return bungalows, err
	}

	return bungalows, nil
}

// GetBungalowByID gets a bungalow by id
func (m *postgresDBRepo) GetBungalowByID(id int) (models.Bungalow, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bungalow models.Bungalow

	query := `
	select 
		id, bungalow_name, created_at, updated_at
	from
		bungalows
	where
		id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&bungalow.ID,
		&bungalow.BungalowName,
		&bungalow.CreatedAt,
		&bungalow.UpdatedAt,
	)

	if err != nil {
		return bungalow, err
	}

	return bungalow, nil
}
