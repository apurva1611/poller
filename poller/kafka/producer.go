package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"poller/model"
	"time"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func Produce(kafkaURL, topic string) {

}

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
}
