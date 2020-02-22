package image

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProcessProducer struct {
	KafkaProducer *kafka.Producer
}

func (p *ProcessProducer) SendProcess(processMsg ProcessMessage, topic string) error {
	value, err := json.Marshal(&processMsg)
	if err != nil {
		return err
	}

	if err := p.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, nil); err != nil {
		return err
	}

	return nil
}
