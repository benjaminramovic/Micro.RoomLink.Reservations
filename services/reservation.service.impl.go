package services

import (
	"context"
	"errors"
	"example/micro-roomlink-reservations/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationServiceImpl struct {
	reservationsCollection *mongo.Collection
	ctx context.Context
}

func NewReservationService(reservationsCollection *mongo.Collection, ctx context.Context) *ReservationServiceImpl {
	return &ReservationServiceImpl{reservationsCollection: reservationsCollection, ctx: ctx}
}

func (rs *ReservationServiceImpl) Create(reservation *models.Reservation) error {
	_, err := rs.reservationsCollection.InsertOne(rs.ctx, reservation)
	
	return err
}
func(rs *ReservationServiceImpl) GetReservation(id string) (*models.Reservation, error) {
	var reservation *models.Reservation
	query := bson.D{bson.E{Key: "id", Value: id}}
	err := rs.reservationsCollection.FindOne(rs.ctx, query).Decode(&reservation)

	return reservation, err
}
func (rs *ReservationServiceImpl) GetAllReservations() ([]*models.Reservation, error) {
	var users []*models.Reservation
	cursor, err := rs.reservationsCollection.Find(rs.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(rs.ctx) {
		var user models.Reservation
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(rs.ctx)

	if len(users) == 0 {
		return nil, errors.New("documents not found")
	}
	return users, nil
}
func (rs *ReservationServiceImpl) Update(reservation *models.Reservation) error {
	return nil
}
func (rs *ReservationServiceImpl) Delete(id string) error {
	return nil
}