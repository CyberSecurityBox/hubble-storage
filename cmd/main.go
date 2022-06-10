package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"hubble-storage/src/core"
	"hubble-storage/src/env"
	"hubble-storage/src/kafka"
	"hubble-storage/src/kafkav2"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var err error
	var receiver kafkav2.ConsumerInterface

	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	handler := core.Core{
		Host:     env.Host,
		Port:     env.Port,
		User:     env.User,
		Password: env.Password,
		Dbname:   env.DbName,
		Log:      log,
	}

	// Consume kafka
	log.Info("Consume started !!!")
	if receiver, err = kafkav2.ConsumeInit(log); err != nil {
		log.WithError(err).Error("Kafka v2 ConsumeInit failed")
		os.Exit(1)
	} else {
		if err = receiver.ReceiveMessage(handler.MainHandler); err != nil {
			log.WithError(err).Error("Kafka v2 ReceiveMessage failed")
			os.Exit(1)
		}
	}
}

func oldConsumer() {
	var err  error
	var consumerInfor *kafka.ConsumerInfor

	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	handler := core.Core{
		Host:     env.Host,
		Port:     env.Port,
		User:     env.User,
		Password: env.Password,
		Dbname:   env.DbName,
		Log:      log,
	}

	// Consumer
	ctx, cancel := context.WithCancel(context.Background())
	chClose := make(chan struct{})
	workersGroup := &sync.WaitGroup{}

	consumerInfor, err = kafka.NewConsumerInfor()
	if err != nil {
		log.WithError(err).Error("Can't create consumer infor")
		return
	}

	consumer, err := kafka.NewConsumer(handler.MainHandler, consumerInfor, log, workersGroup, chClose, ctx)
	if err != nil {
		log.WithError(err).Error("Can't create consumer")
		return
	}
	defer consumer.Close()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Info("Terminating: context cancelled")
		close(chClose)
	case <-sigterm:
		log.Info("Terminating: via signal")
		close(chClose)
	case <-chClose:
		log.Info("Terminating: via close channel")
	}

	cancel()
	workersGroup.Wait()
}


