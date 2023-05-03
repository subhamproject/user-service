package usrmgr

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func SendLogs(msg string) {

	host := GetEnvParam("KAFKA_HOST", "localhost")
	port := GetEnvParam("KAFKA_PORT", "9092")
	topic := GetEnvParam("KAFKA_TOPIC", "demoTopic")

	kafkaServer := host + ":" + port

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServer})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
