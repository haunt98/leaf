package image

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProcessConsumer struct {
	KafkaConsumer *kafka.Consumer
}

func (c *ProcessConsumer) ReceiveProcess() (ProcessMessage, error) {
	msg, err := c.KafkaConsumer.ReadMessage(-1)
	if err != nil {
		return ProcessMessage{}, err
	}

	var processMsg ProcessMessage
	if err := json.Unmarshal(msg.Value, &processMsg); err != nil {
		return ProcessMessage{}, err
	}

	return processMsg, nil
}
