package main

import (
	"ev-payment-service/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	r *gin.Engine
)

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080

}

func InitHttpServer(httpAddress string) {
	r = gin.Default()
	registerHandler()

	if err := r.Run(httpAddress); err != nil {
		logrus.WithField("error", err).Errorf("http server failed to start")
	}
}

func registerHandler() {
	r.GET("/health", handler.GetPaymentHealthCheckHandler)

	v1 := r.Group("/api/v1")
	v1.POST("/provider", handler.CreatePaymentHandler)
}
