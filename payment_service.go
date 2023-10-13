package main

import (
	"context"
	"encoding/json"
	"ev-payment-service/config"
	providerpayment "ev-payment-service/controller/provider"
	userpayment "ev-payment-service/controller/user"
	"ev-payment-service/dao"
	_ "ev-payment-service/docs"
	"ev-payment-service/handler"
	stripeHelper "ev-payment-service/helper"
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

var (
	r *gin.Engine
)

type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type StripeSecret struct {
	StripeKey string `json:"stripekey"`
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
	var documentSecret = retrieveSecretFromSecretManager("EV_DOCUMENTDB")
	if documentSecret.Password != "" {
		mongoHostname = "mongodb://" + documentSecret.Username + ":" + documentSecret.Password + "@docdb-2023-09-23-06-59-34.cluster-cdklkqeyoz4a.ap-southeast-1.docdb.amazonaws.com:27017/?tls=true&tlsCAFile=global-bundle.pem&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false"
	} else {
		mongoHostname = configObj.MongoDBURL
	}
	client := dao.InitDB(mongoHostname)

	var stripeKey string
	var stripeSecret = retrieveStripeKeyFromSecretManager("STRIPE_KEY")
	if stripeSecret.StripeKey != "" {
		stripeKey = stripeSecret.StripeKey
	} else {
		stripeKey = configObj.StripeKey
	}
	userpayment.NewUserController(stripeKey)
	providerpayment.NewProviderPaymentController(stripeKey)
	stripeHelper.StripeKey = stripeKey

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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authentication"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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
			logrus.WithField("decodeSecretManager", database).Error("failed to decode value from secret manager", key)
			return database
		}
	}
	return database
}

func retrieveStripeKeyFromSecretManager(key string) StripeSecret {
	var stripeKey StripeSecret
	secret := os.Getenv(key)
	if secret != "" {
		// Parse the JSON data into the struct
		if err := json.Unmarshal([]byte(secret), &stripeKey); err != nil {
			logrus.WithField("decodeSecretManager", stripeKey).Error("failed to decode value from secret manager", key)
			return stripeKey
		}
	}
	return stripeKey
}

func registerHandler() {

	// use to generate swagger ui
	//	@BasePath	/api/v1
	//	@title		Provider Service API
	//	@version	1.0
	//	@schemes	http
	r.GET("/payment/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/payment/home", handler.GetPaymentHealthCheckHandler)

	v1 := r.Group("/api/v1")
	paymentGroup := v1.Group("/payment")
	paymentGroup.GET("/invoice", handler.GetInvoiceHandler)
	paymentGroup.POST("/invoice", handler.CreateInvoiceHandler)

	paymentGroup.GET("/provider/:provider_id", handler.GetProviderPaymentHandler)
	paymentGroup.POST("/provider", handler.CreateProviderPaymentHandler)
	paymentGroup.PUT("/provider/:provider_id", handler.UpdateProviderPaymentHandler)
	paymentGroup.DELETE("/provider/:provider_id", handler.DeleteProviderPaymentHandler)
	paymentGroup.POST("/provider/completed/:provider_id", handler.CompleteProviderPaymentHandler)

	paymentGroup.GET("/user/:booking_id", handler.GetUserPaymentHandler)
	paymentGroup.POST("/user", handler.CreateUserPaymentHandler)
	paymentGroup.PUT("/user/:booking_id", handler.UpdateUserPaymentHandler)
	paymentGroup.DELETE("/user/:booking_id", handler.DeleteUserPaymentHandler)
	paymentGroup.GET("/user/getAllBooking", handler.GetAllUserPaymentHandler)

	paymentGroup.POST("/user/completed/:booking_id", handler.CompleteUserPaymentHandler)
}
