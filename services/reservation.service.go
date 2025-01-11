package services

import "example/micro-roomlink-reservations/models"

type ReservationService interface {
	Create(reservation *models.Reservation) error
	//GetUserReservations(userId string) ([]models.Reservation, error)
	GetReservation(id string) (*models.Reservation, error)
	GetAllReservations() ([]*models.Reservation, error)
	Update(reservation *models.Reservation) error
	Delete(id string) error
}