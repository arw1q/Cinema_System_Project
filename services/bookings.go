package services

import (
	"Cinema_System_Project/db"
	"Cinema_System_Project/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSession(movieID, hallName string, startTime time.Time, price float64) (*models.Session, error) {
	collection := db.GetCollection("sessions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:        primitive.NewObjectID(),
		MovieID:   objectID,
		HallName:  hallName,
		StartTime: startTime,
		Price:     price,
		CreatedAt: time.Now(),
	}

	_, err = collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func GetSessionsByMovieID(movieID string) ([]models.Session, error) {
	collection := db.GetCollection("sessions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"movie_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func CreateBooking(userID, sessionID, seatNumber string) (*models.Booking, error) {
	collection := db.GetCollection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	sessionObjectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, err
	}

	var existingBooking models.Booking
	err = collection.FindOne(ctx, bson.M{
		"session_id":  sessionObjectID,
		"seat_number": seatNumber,
	}).Decode(&existingBooking)
	if err == nil {
		return nil, errors.New("seat already booked")
	}

	booking := &models.Booking{
		ID:         primitive.NewObjectID(),
		UserID:     userObjectID,
		SessionID:  sessionObjectID,
		SeatNumber: seatNumber,
		CreatedAt:  time.Now(),
	}

	_, err = collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func GetBookingsBySessionID(sessionID string) ([]models.Booking, error) {
	collection := db.GetCollection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"session_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func GetUserBookings(userID string) ([]models.Booking, error) {
	collection := db.GetCollection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func DeleteSession(id string) error {
	collection := db.GetCollection("sessions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func DeleteBooking(id string) error {
	collection := db.GetCollection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func GetAllBookings() ([]models.Booking, error) {
	collection := db.GetCollection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
