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

func main() {
	router := SetupRouter()
	log.Fatal(router.Run(":8080"))

	// get kafka writer using environment variables.
	kafkaURL := "kafka:9092"
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
	kafka.Produce(kafkaURL, produceTopic, minutes)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.GET("/healthcheck", healthCheck)
	return router
}

func healthCheck(c *gin.Context) {
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
	}

	c.JSON(http.StatusOK, "ok")
}
