package example

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"hubble-storage/src/kafka"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func producerExample() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)


	prod, err := kafka.ProducerInit(log)
	if err != nil {
		log.WithError(err).Error("Can't create producer")
		return
	}

	key := `10.40.1.224`
	value := `{"time":"2022-05-26T18:38:46.709177464Z","verdict":"FORWARDED","ethernet":{"source":"76:82:ec:81:8b:77","destination":"82:ec:78:ef:54:ec"},"IP":{"source":"10.40.1.224","destination":"10.148.0.92","ipVersion":"IPv4"},"l4":{"TCP":{"source_port":59380,"destination_port":32001,"flags":{"ACK":true}}},"source":{"ID":3884,"identity":3195,"namespace":"flink-operator-system","labels":["k8s:app=flink","k8s:cluster=flinksessioncluster","k8s:component=taskmanager","k8s:io.cilium.k8s.namespace.labels.control-plane=controller-manager","k8s:io.cilium.k8s.namespace.labels.kubernetes.io/metadata.name=flink-operator-system","k8s:io.cilium.k8s.policy.cluster=gke-secteam-asia-southeast1-b-beyondcorp-uat-cicd","k8s:io.cilium.k8s.policy.serviceaccount=default","k8s:io.kubernetes.pod.namespace=flink-operator-system","k8s:statefulset.kubernetes.io/pod-name=flinksessioncluster-taskmanager-0"],"pod_name":"flinksessioncluster-taskmanager-0","workloads":[{"name":"flinksessioncluster-taskmanager","kind":"StatefulSet"}]},"destination":{"identity":6,"labels":["reserved:remote-node"]},"Type":"L3_L4","node_name":"gke-secteam-asia-southeast1-b-beyondcorp-uat-cicd/gke-beyondcorp-uat-cicd-default-4-8-92745ac7-7c5j","event_type":{"type":4,"sub_type":3},"traffic_direction":"EGRESS","trace_observation_point":"TO_STACK","is_reply":false,"Summary":"TCP Flags: ACK"}`
	_, _, err = prod.SendMessage(key, &value)
	if err != nil {
		log.WithError(err).Error("Can't send message")
		return
	}

	prod.ClearProducer()
}
func consumerExample() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	ctx, cancel := context.WithCancel(context.Background())
	chClose := make(chan struct{})
	workersGroup := &sync.WaitGroup{}

	consumerInfor, err := kafka.NewConsumerInfor()
	if err != nil {
		log.WithError(err).Error("Can't create consumer infor")
		return
	}

	consumer, err := kafka.NewConsumer(test, consumerInfor, log, workersGroup, chClose, ctx)
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
func test(i1 string, i2, i3 []byte) error {
	fmt.Println(i1)
	fmt.Println(i2)
	fmt.Println(i3)
	return nil
}