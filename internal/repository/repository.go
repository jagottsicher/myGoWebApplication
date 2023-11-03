package repository

import (
	"time"

	"github.com/jagottsicher/myGoWebApplication/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertBungalowRestriction(r models.BungalowRestriction) error
	SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowID int) (bool, error)
	SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error)
	GetBungalowByID(id int) (models.Bungalow, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(r models.Reservation) error
	DeleteReservation(id int) error
	UpdateStatusOfReservation(id, status int) error
	AllBungalows() ([]models.Bungalow, error)
	GetRestrictionsForBungalowByDate(bungalowID int, start, end time.Time) ([]models.BungalowRestriction, error)
	InsertBlockForBungalow(id int, startDate time.Time) error
	DeleteBlockByID(id int) error
}
