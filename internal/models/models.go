package models

import "time"

// User is the model of user data
type User struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	Role      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Bungalow is the model of bungalow data
type Bungalow struct {
	ID           int
	BungalowName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Restriction is the model of a restriction
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the model of a reservation
type Reservation struct {
	ID         int
	FullName   string
	Email      string
	Phone      string
	StartDate  time.Time
	EndDate    time.Time
	BungalowID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Bungalow   Bungalow
	Status     int
}

// BungalowRestriction is a model of a bungalow restriction
type BungalowRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	BungalowID    int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Bungalow      Bungalow
	Reservation   Reservation
	Restriction   Restriction
}

// MailData is a model of an e-mail message
type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
