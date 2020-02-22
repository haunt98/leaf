package main

import (
	"leaf/internal/pkg/image"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          image.Group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := kafkaConsumer.SubscribeTopics([]string{image.Topic, "^aRegex.*[Tt]opic"}, nil); err != nil {
		log.Println(err)
		return
	}

	imageProcessConsumer := image.ProcessConsumer{
		KafkaConsumer: kafkaConsumer,
	}

	for {
		processMsg, err := imageProcessConsumer.ReceiveProcess()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("%+v\n", processMsg)
	}
}
