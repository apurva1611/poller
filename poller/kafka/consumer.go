package kafka

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"poller/db"
	"poller/model"
	"strings"

	"github.com/segmentio/kafka-go"
)

func HealthCheck(kafkaURL string) error {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func Consume(kafkaURL, topic, groupID string) {
	reader := getKafkaReader(kafkaURL, topic, groupID)
	defer reader.Close()
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
        log.Info("Consume from id: " + groupID)
		log.Info("CONSUME Topic: %s, Message ID %s", topic, string(m.Key))

		watch := model.WATCH{}
		err = json.Unmarshal(m.Value, &watch)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		messageKey := string(m.Key)

		if strings.HasPrefix(messageKey, "insert") {
			db.InsertWatch(watch)
		} else if strings.HasPrefix(messageKey, "delete") {
			db.DeleteWatch(watch)
		} else if strings.HasPrefix(messageKey, "update") {
			db.UpdateWatch(watch)
		}
	}
}

func TestWeatherMock(kafkaURL, topic, groupID string) {
	reader := getKafkaReader(kafkaURL, topic, groupID)
	defer reader.Close()
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		//log.Print(string(m.Value))

		weather := model.Weather{}
		err = json.Unmarshal(m.Value, &weather)
		if err != nil {
			log.Print(err.Error())
			continue
		}
		log.Info("Topic: " + topic)
		log.Info("Id: " + groupID)
		log.Info("%v : %v", string(m.Key), weather)
	}
}
