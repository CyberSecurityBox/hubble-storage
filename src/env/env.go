package env

import (
	"os"
	"strconv"
)

var (
	Host = os.Getenv("POSTGRES_SERVER")
	Port, _ = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	User = os.Getenv("POSTGRES_USER")
	Password = os.Getenv("POSTGRES_PASS")
	DbName = os.Getenv("POSTGRES_DATABASE")
	KafkaTopic = os.Getenv("KAFKA_TOPIC")
	KafkaUser = os.Getenv("KAFKA_USER")
	KafkaPass = os.Getenv("KAFKA_PASS")
	KafkaBrokers = os.Getenv("KAFKA_BROKERS")
	KafkaGroup = os.Getenv("KAFKA_CONSUME_GROUP")
	KafkaVersion = "2.7.0"
	WebhookUrl = os.Getenv("SLACK_WEBHOOK_URL")
)
