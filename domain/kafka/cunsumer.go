package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"gw-notification/domain/models"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/segmentio/kafka-go"
)

func Consume(wg *sync.WaitGroup, id int, topic string, servers string, groupId string, db *mongo.Database) {
	defer wg.Done()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{servers},
		Topic:       topic,
		GroupID:     groupId,
		StartOffset: kafka.FirstOffset,
		MaxWait:     1 * time.Second,
	})

	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("Consumer %d failed to close reader: %v\n", id, err)
		}
	}()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			if err == context.Canceled {
				log.Printf("Consumer %d shutting down...\n", id)
				return
			}
			if tempErr, ok := err.(interface{ Temporary() bool }); ok && tempErr.Temporary() {
				log.Printf("Consumer %d temporary error: %v\n", id, err)
				continue
			}
			log.Printf("Consumer %d fatal error: %v\n", id, err)
			return
		}

		fmt.Printf("Consumer %d received: %s (partition: %d, offset: %d)\n",
			id, string(msg.Value), msg.Partition, msg.Offset)


		var notification models.Notification
		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			log.Printf("Error unmarshaling notification: %a", err)
		}

		_, err = db.Collection("notifications").InsertOne(context.TODO(), notification)
		if err != nil {
			log.Printf("Error storing message in MongoDB: %s", err)
		}

		if err := r.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("Consumer %d commit error: %v\n", id, err)
		}
	}
}
