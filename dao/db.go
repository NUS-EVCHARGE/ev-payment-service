package dao

import (
	"context"
	"ev-payment-service/dto"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	CreateUserPaymentEntry(userPayment dto.UserPayment) (dto.UserPayment, error)
	GetUserPaymentEntry(id uint) ([]dto.UserPayment, error)
	UpdateUserPaymentEntry(userPayment dto.UserPayment) error
	DeleteUserPaymentEntry(id uint) error
}

var (
	Db Database
)

type dbImpl struct {
	MongoClient *mongo.Client
}

func InitDB(dns string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dns))
	if err != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			logrus.WithField("err", err).Info("error disconnecting from database")
		}
	}
	//defer func() {
	//	if err := client.Disconnect(context.Background()); err != nil {
	//		panic(err)
	//	}
	//}()

	// Send a ping to confirm a successful connection
	if err = client.Ping(context.Background(), nil); err != nil {
		logrus.WithField("err", err).Info("error pinging database")
	} else {
		logrus.WithField("success", "connected to database").Info("connected to database")
	}

	Db = NewDatabase(client)
}

func NewDatabase(client *mongo.Client) Database {
	return &dbImpl{
		MongoClient: client,
	}
}
