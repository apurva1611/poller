package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"poller/db"
	"poller/kafka"
	"strconv"

	"github.com/gin-gonic/gin"
)

const kafkaURL = "kafka:9092"

func main() {
	// get kafka writer using environment variables.
	consumeTopic := "watch"
	consumeGroup := "watch-group"

	produceTopic := "weather"
	minutesStr := os.Getenv("minutes")

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		// default
		minutes = 1
	}

	db.Init()
	defer db.CloseDB()

	go kafka.Consume(kafkaURL, consumeTopic, consumeGroup)

	fmt.Println("now to producer")
	go kafka.Produce(kafkaURL, produceTopic, minutes)

	router := SetupRouter()
	log.Fatal(router.Run(":8080"))
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.GET("/healthcheck", healthCheck)
	return router
}

func healthCheck(c *gin.Context) {
<<<<<<< HEAD
	kafkaURL := "kafka:9092"
	err := db.HealthCheck()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "db health check failed.")
		os.Exit(3)
	}

	err = kafka.KafkaHealthCheck(kafkaURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "kafka health check failed.")
		os.Exit(4)
=======
	err := db.HealthCheck()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "db health check failed.")
		os.Exit(5)
	}

	err = kafka.HealthCheck(kafkaURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "kafka health check failed.")
		os.Exit(6)
>>>>>>> 417936edb0b368ff2e84e3bed023ab49d51ed6e6
	}

	c.JSON(http.StatusOK, "ok")
}
