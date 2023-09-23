package main

import (
	"context"
	"encoding/json"
	"ev-payment-service/config"
	userpayment "ev-payment-service/controller/user"
	"ev-payment-service/dao"
	"ev-payment-service/handler"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

var (
	r *gin.Engine
)

type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

	var mongoHostname string
	var documentSecret = retrieveSecretFromSecretManager("EV_DOCUMENTDV")
	if documentSecret.Password != "" {
		mongoHostname = "mongodb://" + documentSecret.Username + ":" + documentSecret.Password + "@docdb-2023-09-23-06-59-34.cluster-cdklkqeyoz4a.ap-southeast-1.docdb.amazonaws.com:27017/?tls=true&tlsCAFile=global-bundle.pem&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false"
	} else {
		mongoHostname = configObj.MongoDBURL
	}
	client := dao.InitDB(mongoHostname)

	var stripeKey string
	var stripeSecret = retrieveSecretFromSecretManager("STRIPE_KEY")
	if stripeSecret.Password != "" {
		stripeKey = stripeSecret.Password
	} else {
		stripeKey = configObj.StripeKey
	}
	userpayment.NewUserController(stripeKey)

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

func retrieveSecretFromSecretManager(key string) DatabaseSecret {
	var database DatabaseSecret
	secret := os.Getenv(key)
	if secret != "" {
		// Parse the JSON data into the struct
		if err := json.Unmarshal([]byte(secret), &database); err != nil {
			logrus.WithField("decodeSecretManager", database).Error("failed to decode value from secret manager")
			return database
		}
	}
	return database
}

func registerHandler() {
	r.GET("/payment/home", handler.GetPaymentHealthCheckHandler)

	v1 := r.Group("/api/v1")
	paymentGroup := v1.Group("/payment")
	paymentGroup.POST("/provider", handler.CreatePaymentHandler)

	paymentGroup.GET("/user/:booking_id", handler.GetUserPaymentHandler)
	paymentGroup.POST("/user", handler.CreateUserPaymentHandler)
	paymentGroup.PUT("/user/:booking_id", handler.UpdateUserPaymentHandler)
	paymentGroup.DELETE("/user/:booking_id", handler.DeleteUserPaymentHandler)
}
