package dao

import (
	"context"
	"ev-payment-service/dto"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	CreateUserPaymentEntry(userPayment *dto.UserPayment) error
	GetUserPaymentEntry(bookingId uint) ([]dto.UserPayment, error)
	GetUserPaymentByProviderId(providerId uint) ([]dto.UserPayment, error)
	UpdateUserPaymentEntry(userPayment *dto.UserPayment) error
	DeleteUserPaymentEntry(id uint) error

	CreateProviderPaymentEntry(providerPayment *dto.ProviderPayment) error
	GetProviderPaymentEntry(providerId uint, billingPeriod dto.ProviderBillingPeriod) ([]dto.ProviderPayment, error)
	UpdateProviderPaymentEntry(providerPayment *dto.ProviderPayment) error
	DeleteProviderPaymentEntry(id uint, billingPeriod dto.ProviderBillingPeriod) error

	Close()
}

var (
	Db Database
)

type DbImpl struct {
	MongoClient *mongo.Client
}

func (db *DbImpl) Close() {
	err := db.MongoClient.Disconnect(context.Background())
	if err != nil {
		logrus.WithField("err", err).Info("error disconnecting from database")
	}
}

func InitDB(dns string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dns))
	if err != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			logrus.WithField("err", err).Info("error disconnecting from database")
		}
	}

	// Send a ping to confirm a successful connection
	if err = client.Ping(context.Background(), nil); err != nil {
		logrus.WithField("err", err).Info("error pinging database")
	} else {
		logrus.WithField("success", "connected to database").Info("connected to database")
	}

	Db = NewDatabase(client)
}

func NewDatabase(client *mongo.Client) Database {
	return &DbImpl{
		MongoClient: client,
	}
}
