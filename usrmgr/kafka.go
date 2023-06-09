package usrmgr

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/subhamproject/user-service/utils"
)

var (
	kafkaWriter *kafka.Writer
	topic       string
)

func SendLogs(val string) {
	fmt.Println("writing event to kafka", val)
	msg := kafka.Message{
		Partition: int(kafka.PatternTypeAny),
		Value:     []byte(val),
	}

	err := kafkaWriter.WriteMessages(context.TODO(), msg)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("sucessfully sent message to kafka, event is: ", val)
	}
}

func getSslKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	clientCertFile := utils.GetEnvParam("KAFKA_CLIENT_CERT", "/home/om/go/src/github.com/subhamproject/devops-demo/certs/kafka.user.cert")
	clientKeyFile := utils.GetEnvParam("KAFKA_CLIENT_KEY", "/home/om/go/src/github.com/subhamproject/devops-demo/certs/kafka.user.key")
	// caCertFile := GetEnvParam("KAFKA_CA_CERT", "/home/om/go/src/github.com/subhamproject/devops-demo/certs/kafka.user.pem")
	servers := strings.Split(kafkaURL, ",")
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       cfg,
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  servers,
		Topic:    topic,
		Balancer: &kafka.Hash{},
		Dialer:   dialer,
	})
	w.AllowAutoTopicCreation = true

	return w
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	fmt.Println("initializing non ssl kafka writer")
	servers := strings.Split(kafkaURL, ",")

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  servers,
		Topic:    topic,
		Balancer: &kafka.Hash{},
	})
	w.AllowAutoTopicCreation = true

	// w := &kafka.Writer{
	// 	Addr:                   kafka.TCP(kafkaURL),
	// 	Topic:                  topic,
	// 	AllowAutoTopicCreation: true,
	// }

	return w
}

func InitKafka() {
	devMode := utils.GetEnvBoolParam("DEV_MODE", true)
	fmt.Println("initializing kafka connection. devmode: ", devMode)
	kafkaURL := utils.GetEnvParam("KAFKA_SERVERS", "localhost:9092")
	servers := strings.Split(kafkaURL, ",")
	fmt.Println("kafka servers: ", servers)

	// get kafka writer using environment variables.
	topic = utils.GetEnvParam("KAFKA_TOPIC", "demoTopic")
	if devMode {
		kafkaWriter = getKafkaWriter(kafkaURL, topic)
	} else {
		kafkaWriter = getSslKafkaWriter(kafkaURL, topic)
	}

	fmt.Println("init kafka writer - ", kafkaWriter)

	if err := pingKafka(); err != nil {
		log.Fatal("kafka ping error: ", err)
	}
}

func CloseKafka() {
	if err := kafkaWriter.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func pingKafka() error {

	messages := []kafka.Message{
		{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	}

	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = kafkaWriter.WriteMessages(ctx, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Fatalf("unexpected error %v", err)
		}
	}
	return err
}
