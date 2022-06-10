package kafka

import (
	"context"
	"fmt"
	"hubble-storage/src/env"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type ConsumerInfor struct {
	KafkaBrokers       string
	KafkaUser          string
	KafkaPassword      string
	KafkaVersion       string
	KafkaConsumerTopic string
	KafkaConsumerGroup string
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	fnHandleData func(string, []byte, []byte) error
	log          logrus.FieldLogger
	ready        chan bool
	group        string
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		if err := consumer.fnHandleData(message.Topic, message.Key, message.Value); err != nil {
			consumer.log.WithError(err).Error("Could not handle data")
		} else {
			//consumer.log.Infof("Partition: %v === Group: %v", message.Partition, consumer.group)
			session.MarkMessage(message, "")
		}
	}

	return nil
}

func NewConsumer(fnHandleData func(string, []byte, []byte) error, consumerInfor *ConsumerInfor, log logrus.FieldLogger, wg *sync.WaitGroup, chClose chan struct{}, ctx context.Context) (sarama.ConsumerGroup, error) {
	var scramClient = XDGSCRAMClient{
		Client:             nil,
		ClientConversation: nil,
		HashGeneratorFcn:   SHA512,
	}

	version, err := sarama.ParseKafkaVersion(consumerInfor.KafkaVersion)
	if err != nil {
		log.WithError(err).Error("Can't parse kafka version")
		return nil, err
	}

	configConsumer := sarama.NewConfig()

	configConsumer.Version = version
	configConsumer.Consumer.Group.Rebalance.Strategy = MyBalanceStrategy
	configConsumer.Consumer.Offsets.Initial = sarama.OffsetOldest
	configConsumer.Consumer.Return.Errors = true
	configConsumer.Metadata.Full = true //default

	configConsumer.Net.SASL.Enable = true
	//configConsumer.Net.SASL.Handshake = true
	configConsumer.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	configConsumer.Net.SASL.User = consumerInfor.KafkaUser
	configConsumer.Net.SASL.Password = consumerInfor.KafkaPassword
	configConsumer.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scramClient }

	client, err := sarama.NewConsumerGroup(strings.Split(consumerInfor.KafkaBrokers, ","), consumerInfor.KafkaConsumerGroup, configConsumer)
	if err != nil {
		log.WithError(err).Error("Can't create consumer group")
		return nil, err
	}

	groupConsumer := Consumer{
		fnHandleData: fnHandleData,
		log:          log,
		ready:        make(chan bool),
		group:        consumerInfor.KafkaConsumerGroup,
	}

	//log error
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-chClose:
				return
			case err := <-client.Errors():
				if err != nil {
					log.WithError(err).Error("Consumer get a error")
				}
			}
		}
	}()

	//consume
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(consumerInfor.KafkaConsumerTopic, ","), &groupConsumer); err != nil {
				log.WithError(err).Error("Can't consume topics")
				close(chClose)
				return
			}

			if ctx.Err() != nil {
				log.WithError(ctx.Err()).Error("Context was cancelled")
				close(chClose)
				return
			}

			groupConsumer.ready = make(chan bool)
		}
	}()

	select {
	case <-groupConsumer.ready:
		log.Info("Sarama consumer up and running!...")
	case <-time.After(time.Minute):
		log.Error("Consume topics timeout")
		client.Close()
		return nil, fmt.Errorf("Consume topics timeout.")
	}

	return client, nil
}

func NewConsumerInfor() (*ConsumerInfor, error) {
	kafkaBrokers := env.KafkaBrokers
	if kafkaBrokers == "" {
		return nil, fmt.Errorf("Missing kafka brokers in env")
	}

	kafkaUser := env.KafkaUser
	if kafkaUser == "" {
		return nil, fmt.Errorf("Missing kafka user in env")
	}

	kafkaPassword := env.KafkaPass
	if kafkaPassword == "" {
		return nil, fmt.Errorf("Missing kafka password in env")
	}

	kafkaConsumerTopic := env.KafkaTopic
	if kafkaConsumerTopic == "" {
		return nil, fmt.Errorf("Missing kafka consumer topic in env")
	}

	kafkaConsumerGroup := env.KafkaGroup
	if kafkaConsumerGroup == "" {
		return nil, fmt.Errorf("Missing kafka consumer group in env")
	}

	return &ConsumerInfor{
		KafkaBrokers:       kafkaBrokers,
		KafkaUser:          kafkaUser,
		KafkaPassword:      kafkaPassword,
		KafkaConsumerTopic: kafkaConsumerTopic,
		KafkaConsumerGroup: kafkaConsumerGroup,
		KafkaVersion:       "2.7.0",
	}, nil
}
