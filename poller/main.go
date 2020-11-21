package main

import (
	"net/http"
	"os"
	"poller/db"
	"poller/kafka"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const kafkaURL = "app-prereq-kafka:9092"

func main() {
	db.Init()
	defer db.CloseDB()

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

	go kafka.Consume(kafkaURL, consumeTopic, consumeGroup)

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
	//kafkaURL := "app-prereq-kafka.monitoring:9092"
	err := db.HealthCheck()
	if err != nil {
		log.Error("DB HEALTHCHECK: %s", err.Error())
		log.Printf("DB HEALTHCHECK %s", err.Error())
		c.JSON(http.StatusInternalServerError, "db health check failed.")
		os.Exit(3)
	}

	err = kafka.HealthCheck(kafkaURL)

	if err != nil {
		log.Error("KAFKA HEALTHCHECK %s", err.Error())
		c.JSON(http.StatusInternalServerError, "kafka health check failed.")
		os.Exit(4)
	}

	log.Info("HEALTHCHECK SUCCESSFUL")

	c.JSON(http.StatusOK, "ok")
}
