package models

import "time"

// Reservation contains reservation data
type Reservation struct {
	Name  string
	Email string
	Phone string
}

// User is the model of user data
type User struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	Role      int
	CreatedAt time.Time
	UdpadedAt time.Time
}

// Bungalow is the model of bungalow data
type Bungalow struct {
	ID           int
	BungalowName string
	CreatedAt    time.Time
	UdpadedAt    time.Time
}

// Restrictions is the model of a restriction
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UdpadedAt       time.Time
}

// Reservations is the model of a reservation
type Reservations struct {
	ID         int
	FullName   string
	Email      string
	Phone      string
	StartDate  time.Time
	EndDate    time.Time
	BungalowID int
	CreatedAt  time.Time
	UdpadedAt  time.Time
	Bungalow   Bungalow
	Processed  int
}

// BungalowRestrictions is a model of a bungalow restriction
type BungalowRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	BungalowID    int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Bungalow      Bungalow
	Reservation   Reservations
	Restriction   Restrictions
}
