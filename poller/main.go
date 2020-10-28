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
	// get kafka writer using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
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
	v1.GET("/healthcheck", healthcheck)
	return router
}

func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
