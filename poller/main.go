package main

import (
	"os"
	"poller/db"
	"poller/kafka"
	"strconv"
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
		minutes = 10
	}

	db.Init()
	defer db.CloseDB()

	go kafka.Consume(kafkaURL, consumeTopic, consumeGroup)
	kafka.Produce(kafkaURL, produceTopic, minutes)
}
