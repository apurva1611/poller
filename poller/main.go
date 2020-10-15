package main

import (
	"os"
	"poller/db"
	"poller/kafka"
)

func main() {
	// get kafka writer using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	consumeTopic := "watch"
	consumeGroup := "watch-group"

	produceTopic := "weather"

	db.Init()
	defer db.CloseDB()

	go kafka.Consume(kafkaURL, consumeTopic, consumeGroup)
	kafka.Produce(kafkaURL, produceTopic)
}
