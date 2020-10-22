package main

import (
	"fmt"
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
		minutes = 5
	}

	db.Init()
	defer db.CloseDB()

	go kafka.Consume(kafkaURL, consumeTopic, consumeGroup)

	fmt.Println("now to producer")
	kafka.Produce(kafkaURL, produceTopic, minutes)
}
