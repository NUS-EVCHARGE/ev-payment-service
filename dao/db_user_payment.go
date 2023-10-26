package dao

import (
	"context"
	"ev-payment-service/dto"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func (db *DbImpl) GetUserPaymentEntry(id uint) ([]dto.UserPayment, error) {

	userCollection := db.MongoClient.Database("ev").Collection("user_payment")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filterCursor, err := userCollection.Find(ctx, bson.M{"bookingId": id})
	if err != nil {
		logrus.WithField("err", err).Info("error getting user payment")
		return nil, err
	}

	var userPaymentFiltered []dto.UserPayment
	if err = filterCursor.All(ctx, &userPaymentFiltered); err != nil {
		logrus.WithField("err", err).Info("error getting user payment")
		return nil, err
	} else {
		logrus.WithField("success", userPaymentFiltered).Info("user payment retrieved")
		return userPaymentFiltered, nil
	}
}

func (db *DbImpl) GetUserPaymentByProviderId(providerId uint) ([]dto.UserPayment, error) {

	userCollection := db.MongoClient.Database("ev").Collection("user_payment")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filterCursor, err := userCollection.Find(ctx, bson.M{"providerId": providerId})
	if err != nil {
		logrus.WithField("err", err).Info("error getting user payment")
		return nil, err
	}

	var userPaymentFiltered []dto.UserPayment
	if err = filterCursor.All(ctx, &userPaymentFiltered); err != nil {
		logrus.WithField("err", err).Info("error getting user payment")
		return nil, err
	} else {
		logrus.WithField("success", userPaymentFiltered).Info("user payment retrieved")
		return userPaymentFiltered, nil
	}
}

func (db *DbImpl) CreateUserPaymentEntry(userPayment *dto.UserPayment) error {

	userCollection := db.MongoClient.Database("ev").Collection("user_payment")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	items, err := Db.GetUserPaymentEntry(userPayment.BookingId)
	if items != nil {
		logrus.WithField("items", items).Info("user payment for booking id already exists")
		return fmt.Errorf("user payment booking id already exists")
	}
	createUserPaymentResult, err := userCollection.InsertOne(ctx, &userPayment)
	if err != nil {
		logrus.WithField("err", err).Info("error inserting user payment")
		return err
	} else {
		logrus.WithField("success", createUserPaymentResult.InsertedID).Info("user payment inserted")
	}
	return nil
}

func (db *DbImpl) UpdateUserPaymentEntry(userPayment *dto.UserPayment) error {
	collection := db.MongoClient.Database("ev").Collection("user_payment")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"bookingId": userPayment.BookingId}
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": userPayment})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (db *DbImpl) DeleteUserPaymentEntry(id uint) error {
	collection := db.MongoClient.Database("ev").Collection("user_payment")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"bookingId": id}
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
