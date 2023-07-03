package repository

import "github.com/jagottsicher/myGoWebApplication/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
