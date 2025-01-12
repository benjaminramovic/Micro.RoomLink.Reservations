package services

import (
	"context"
	"errors"
	"example/micro-roomlink-reservations/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationServiceImpl struct {
	reservationsCollection *mongo.Collection
	ctx context.Context
}

func NewReservationService(reservationsCollection *mongo.Collection, ctx context.Context) *ReservationServiceImpl {
	return &ReservationServiceImpl{reservationsCollection: reservationsCollection, ctx: ctx}
}

//CRUD
func (rs *ReservationServiceImpl) Create(reservation *models.Reservation) error {
	if reservation.Id.IsZero() {
        reservation.Id = primitive.NewObjectID()
    }
	_, err := rs.reservationsCollection.InsertOne(rs.ctx, reservation)
	
	return err
}
func (rs *ReservationServiceImpl) GetReservation(id string) (*models.Reservation, error) {
	// Konverzija stringa u ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil,errors.New("invalid ID format")
	}
	// Kreiranje filtera prema _id polju
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	// Traženje dokumenta
	var reservation models.Reservation
	err = rs.reservationsCollection.FindOne(rs.ctx, filter).Decode(&reservation)
	if err != nil {
		return nil, err // Vraća mongo.ErrNoDocuments ako nema rezultata
	}

	return &reservation, nil
}
func (rs *ReservationServiceImpl) GetAllReservations() ([]*models.Reservation, error) {
	var reservations []*models.Reservation
	cursor, err := rs.reservationsCollection.Find(rs.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(rs.ctx) {
		var reservation models.Reservation
		err := cursor.Decode(&reservation)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(rs.ctx)

	if len(reservations) == 0 {
		return nil, errors.New("documents not found")
	}
	return reservations, nil
}
func (rs *ReservationServiceImpl) Update(id string, reservation *models.Reservation) error {
	// Filter za pronalaženje dokumenta po _id (kao string)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	// Ažuriranje polja
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			{Key: "check_in", Value: reservation.CheckIn},
			{Key: "check_out", Value: reservation.CheckOut},
			{Key: "status", Value: reservation.Status},
			{Key: "total_price", Value: reservation.TotalPrice},
		}},
	}

	// Izvršavanje ažuriranja
	result, err := rs.reservationsCollection.UpdateOne(rs.ctx, filter, update)
	if err != nil {
		return err // Vraća grešku ako ažuriranje nije uspelo
	}

	// Provera da li je dokument pronađen
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}

	return nil // Uspešno ažuriranje
}

func (rs *ReservationServiceImpl) Delete(id string) error {
	// Kreiranje filtera prema _id polju
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	// Brisanje dokumenta
	result, err := rs.reservationsCollection.DeleteOne(rs.ctx, filter)
	if err != nil {
		return err // Greška pri izvršavanju zahteva
	}

	// Provera da li je dokument pronađen i obrisan
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}

	return nil // Uspešno brisanje
}
func (rs *ReservationServiceImpl) GetGuestReservations(guestId int) ([]*models.Reservation, error) {
	var reservations []*models.Reservation
	filter := bson.D{primitive.E{Key: "guest_id", Value: guestId}}
	cursor, err := rs.reservationsCollection.Find(rs.ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(rs.ctx) {
		var reservation models.Reservation
		err := cursor.Decode(&reservation)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(rs.ctx)

	if len(reservations) == 0 {
		return nil, errors.New("documents not found")
	}
	return reservations, nil
}
