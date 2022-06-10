package kafkav2

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type consumer struct {
	context context.Context
	dialer  *kafka.Dialer
	log     *logrus.Logger
}
type ConsumerInterface interface {
	ReceiveMessage(fnHandleData func(string, []byte, []byte) error) error
}

type ProducerInterface interface {
	Key(uuid string, name string, scanType string) []byte
	SendMessage(key []byte, value []byte) error
	close()
}

