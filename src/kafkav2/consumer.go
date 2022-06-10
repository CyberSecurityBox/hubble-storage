package kafkav2

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
	"hubble-storage/src/env"
	"time"
)

func ConsumeInit(log *logrus.Logger) (ConsumerInterface, error) {
	/*
	 - Authentication with scram.SHA512 mechanism
	 - Return dialer after authenticated
	*/

	var (
		ctx = context.Background()
		mechanism sasl.Mechanism
		dialer *kafka.Dialer
		err error
	)

	if mechanism, err = scram.Mechanism(scram.SHA512, env.KafkaUser, env.KafkaPass); err != nil {
		return nil, err
	}

	dialer = &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	return &consumer{context: ctx, dialer: dialer, log: log}, nil
}
func (c *consumer) ReceiveMessage(fnHandleData func(string, []byte, []byte) error) error {
	var reader *kafka.Reader

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{env.KafkaBrokers},
		Topic:   env.KafkaTopic,
		GroupID: env.KafkaGroup,
		Dialer:  c.dialer,
	})

	defer reader.Close()

	for {
		var message kafka.Message
		var	err error

		if message, err = reader.FetchMessage(c.context); err != nil {
			c.log.WithError(err).Error("Fetching message")
			return err
		}

		if err = reader.CommitMessages(c.context, message); err != nil {
			c.log.WithError(err).Error("Commit messages")
			return err
		}

		// Handle
		err = fnHandleData(message.Topic, message.Key, message.Value)
		if err != nil {
			c.log.WithError(err).Error("fnHandleData error")
			return err
		}
	}
}

