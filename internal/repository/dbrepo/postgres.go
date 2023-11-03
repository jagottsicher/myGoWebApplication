package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jagottsicher/myGoWebApplication/internal/models"
	"golang.org/x/crypto/bcrypt"
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
			and $2 <= end_date and $3 >= start_date;
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
	defer rows.Close()

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

// GetUserByID returns user data by id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, full_name, email, password, role, created_at, updated_at
	from users where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FullName,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates basic user data in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update users set full_name = $1, email = $2, role = $3, updated_at = $4
`
	_, err := m.DB.ExecContext(ctx, query,
		u.FullName,
		u.Email,
		u.Role,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user by data
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var passwordHash string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email =$1", email)

	err := row.Scan(&id, &passwordHash)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("wrong password")
	} else if err != nil {
		return 0, "", err
	}

	return id, passwordHash, nil
}

// AllReservations builds and returns a slice of all reservations from the database
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.full_name, r.email, r.phone, r.start_date, 
		r.end_date, r.bungalow_id, r.created_at, r.updated_at, r.status,
		b.id, b.bungalow_name
		from reservations r
		left join bungalows b on (r.bungalow_id = b.id)
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.BungalowID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Status,
			&i.Bungalow.ID,
			&i.Bungalow.BungalowName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations builds and returns a slice of all new reservations from the database
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.full_name, r.email, r.phone, r.start_date, 
		r.end_date, r.bungalow_id, r.created_at, r.updated_at, r.status,
		b.id, b.bungalow_name
		from reservations r
		left join bungalows b on (r.bungalow_id = b.id)
		where status = 0
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.BungalowID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Status,
			&i.Bungalow.ID,
			&i.Bungalow.BungalowName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// GetReservationByID returns a reservation by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
		select r.id, r.full_name, r.email, r.phone, r.start_date, 
		r.end_date, r.bungalow_id, r.created_at, r.updated_at, r.status,
		b.id, b.bungalow_name
		from reservations r
		left join bungalows b on (r.bungalow_id = b.id)
		where r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FullName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.BungalowID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Status,
		&res.Bungalow.ID,
		&res.Bungalow.BungalowName,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateReservation updates the data of a reservation in the database
func (m *postgresDBRepo) UpdateReservation(r models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update reservations set full_name = $1, email = $2, phone = $3, updated_at = $4
		where id = $5
`
	_, err := m.DB.ExecContext(ctx, query,
		r.FullName,
		r.Email,
		r.Phone,
		time.Now(),
		r.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation by id deletes an entry of a reservation dron the database
func (m *postgresDBRepo) DeleteReservation(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		delete from reservations
		where id = $1
`
	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

// UpdateStatusOfReservation by id updates the status of a reservation
func (m *postgresDBRepo) UpdateStatusOfReservation(id, status int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update reservations set status = $1
		where id = $2
`
	_, err := m.DB.ExecContext(ctx, query, status, id)

	if err != nil {
		return err
	}

	return nil
}

// AllBungalows returns a slice of all bungalows
func (m *postgresDBRepo) AllBungalows() ([]models.Bungalow, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bungalows []models.Bungalow

	query := `select id, bungalow_name, created_at, updated_at from bungalows order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return bungalows, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Bungalow
		err := rows.Scan(
			&b.ID,
			&b.BungalowName,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return bungalows, err
		}
		bungalows = append(bungalows, b)
	}

	if err = rows.Err(); err != nil {
		return bungalows, err
	}

	return bungalows, nil
}

// GetRestrictionsForBungalowByDate returns restrictions for a bungalow by date range
func (m *postgresDBRepo) GetRestrictionsForBungalowByDate(bungalowID int, start, end time.Time) ([]models.BungalowRestriction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.BungalowRestriction

	query := `
		select id, coalesce(reservation_id, 0), restriction_id, bungalow_id, start_date, end_date
		from bungalow_restrictions where $1 < end_date and $2 >= start_date
		and bungalow_id = $3
	`
	rows, err := m.DB.QueryContext(ctx, query, start, end, bungalowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.BungalowRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.BungalowID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil

}

// InsertBlockForBungalow inserts a bungalow restriction by bungalow id for a specific day
func (m *postgresDBRepo) InsertBlockForBungalow(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into bungalow_restrictions (start_date, end_date, bungalow_id, restriction_id,
			created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteBlockByID deletes a bungalow restriction by id
func (m *postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from bungalow_restrictions where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
