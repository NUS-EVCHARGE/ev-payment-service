package main

import (
	"context"
	"ev-payment-service/config"
	userpayment "ev-payment-service/controller/user"
	"ev-payment-service/dao"
	"ev-payment-service/handler"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	r *gin.Engine
)

func main() {

	var (
		configFile string
	)

	flag.StringVar(&configFile, "config", "config.yaml", "configuration file of this service")
	flag.Parse()

	// init configurations
	configObj, err := config.ParseConfig(configFile)
	if err != nil {
		logrus.WithField("error", err).WithField("filename", configFile).Error("failed to init configurations")
		return
	}

	// init db
	var mongoHostname string = configObj.MongoDBURL

	userpayment.NewUserController(configObj.StripeKey)
	client := dao.InitDB(mongoHostname)

	//defer disconnect from database
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			logrus.WithField("err", err).Info("error disconnecting from database")
		}
	}(client, context.Background())

	InitHttpServer(configObj.HttpAddress)

}

func InitHttpServer(httpAddress string) {
	r = gin.Default()
	registerHandler()

	if err := r.Run(httpAddress); err != nil {
		logrus.WithField("error", err).Errorf("http server failed to start")
	}
}

func registerHandler() {
	r.GET("/payment/home", handler.GetPaymentHealthCheckHandler)

	v1 := r.Group("/api/v1")
	v1.POST("/provider", handler.CreatePaymentHandler)

	v1.GET("/userpayment/:booking_id", handler.GetUserPaymentHandler)
	v1.POST("/userpayment", handler.CreateUserPaymentHandler)
	v1.PUT("/userpayment/:booking_id", handler.UpdateUserPaymentHandler)
	v1.DELETE("/userpayment/:booking_id", handler.DeleteUserPaymentHandler)
}
