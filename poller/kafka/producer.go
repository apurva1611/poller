package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"poller/db"
	"poller/model"
	"time"

	"github.com/segmentio/kafka-go"
)

func KafkaHealthCheck(kafkaURL string) error {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func Produce(kafkaURL, topic string, minutes int) {
	fmt.Println("create new writer")
	writer := newKafkaWriter(kafkaURL, topic)
	log.Print("new kafka writer in created")
	for {
		time.Sleep(time.Duration(minutes) * time.Minute)
		zipCodesSet := db.GetAllZipCodes()
		for _, zipCode := range zipCodesSet {
			weather := GetWeatherData(zipCode)

			weatherTopicData := model.WeatherTopicData{}
			weatherTopicData.Zipcode = zipCode
			weatherTopicData.WeatherData = *weather
			weatherTopicData.Watchs = db.GetAllWatchesByZipcode(zipCode)
			weatherTopicDataJSON, _ := json.Marshal(weatherTopicData)
			fmt.Println("getting watches as per zipcode in producer: " + string(weatherTopicDataJSON))
			msg := kafka.Message{
				Key:   []byte(zipCode),
				Value: weatherTopicDataJSON,
			}

			err := writer.WriteMessages(context.Background(), msg)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("producer success")
		}
	}
}

/*
func TestWatchMock(kafkaURL, topic string) {
	writer := newKafkaWriter(kafkaURL, topic)
	for i := 0; ; i++ {
		watchID := uuid.New().String()
		var alerts []model.ALERT
		alerts = append(alerts, model.ALERT{
			ID:           uuid.New().String(),
			WatchId:      watchID,
			FieldType:    "temp_min",
			Operator:     "gt",
			Value:        50,
			AlertCreated: "2006-01-02 03:04:05",
			AlertUpdated: "2006-01-02 03:04:05",
		})

		watch := model.WATCH{
			ID:           watchID,
			UserId:       uuid.New().String(),
			Zipcode:      "02115",
			Alerts:       alerts,
			WatchCreated: "2006-01-02 03:04:05",
			WatchUpdated: "2006-01-02 03:04:05",
		}

		watchJson, _ := json.Marshal(watch)

		msg := kafka.Message{
			Key:   []byte("insert" + watch.Alerts[0].ID),
			Value: watchJson,
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
} */
