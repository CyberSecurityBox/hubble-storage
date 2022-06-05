package example

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "pass"
	dbname   = "db"
)

func pgsqlExample() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	// insert
	// hardcoded
	insertStmt := `insert into "flows"(flow) values('{"time":"2022-05-26T18:38:46.709177464Z","verdict":"FORWARDED","ethernet":{"source":"76:82:ec:81:8b:77","destination":"82:ec:78:ef:54:ec"},"IP":{"source":"10.40.1.224","destination":"10.148.0.92","ipVersion":"IPv4"},"l4":{"TCP":{"source_port":59380,"destination_port":32001,"flags":{"ACK":true}}},"source":{"ID":3884,"identity":3195,"namespace":"flink-operator-system","labels":["k8s:app=flink","k8s:cluster=flinksessioncluster","k8s:component=taskmanager","k8s:io.cilium.k8s.namespace.labels.control-plane=controller-manager","k8s:io.cilium.k8s.namespace.labels.kubernetes.io/metadata.name=flink-operator-system","k8s:io.cilium.k8s.policy.cluster=gke-secteam-asia-southeast1-b-beyondcorp-uat-cicd","k8s:io.cilium.k8s.policy.serviceaccount=default","k8s:io.kubernetes.pod.namespace=flink-operator-system","k8s:statefulset.kubernetes.io/pod-name=flinksessioncluster-taskmanager-0"],"pod_name":"flinksessioncluster-taskmanager-0","workloads":[{"name":"flinksessioncluster-taskmanager","kind":"StatefulSet"}]},"destination":{"identity":6,"labels":["reserved:remote-node"]},"Type":"L3_L4","node_name":"gke-secteam-asia-southeast1-b-beyondcorp-uat-cicd/gke-beyondcorp-uat-cicd-default-4-8-92745ac7-7c5j","event_type":{"type":4,"sub_type":3},"traffic_direction":"EGRESS","trace_observation_point":"TO_STACK","is_reply":false,"Summary":"TCP Flags: ACK"}')`
	_, e := db.Exec(insertStmt)
	CheckError(e)

	// dynamic
	insertDynStmt := `insert into "flows"(flow) values($1)`
	_, e = db.Exec(insertDynStmt, `{"time":"2022-05-26T18:38:46.709177464Z","verdict":"FORWARDED","ethernet":{"source":"76:82:ec:81:8b:77","destination":"82:ec:78:ef:54:ec"},"IP":{"source":"10.40.1.224","destination":"10.148.0.92","ipVersion":"IPv4"},"l4":{"TCP":{"source_port":59380,"destination_port":32001,"flags":{"ACK":true}}},"source":{"ID":3884,"identity":3195,"namespace":"flink-operator-system","labels":["k8s:app=flink","k8s:cluster=flinksessioncluster","k8s:component=taskmanager","k8s:io.cilium.k8s.namespace.labels.control-plane=controller-manager","k8s:io.cilium.k8s.namespace.labels.kubernetes.io/metadata.name=flink-operator-system","k8s:io.cilium.k8s.policy.cluster=gke-secteam-asia-southeast1-b-beyondcorp-uat-cicd","k8s:io.cilium.k8s.policy.serviceaccount=default","k8s:io.kubernetes.pod.namespace=flink-operator-system","k8s:statefulset.kubernetes.io/pod-name=flinksessioncluster-taskmanager-0"],"pod_name":"flinksessioncluster-taskmanager-0","workloads":[{"name":"flinksessioncluster-taskmanager","kind":"StatefulSet"}]},"destination":{"identity":6,"labels":["reserved:remote-node"]},"Type":"L3_L4","node_name":"gke-secteam-asia-southeast1-b-beyondcorp-uat-cicd/gke-beyondcorp-uat-cicd-default-4-8-92745ac7-7c5j","event_type":{"type":4,"sub_type":3},"traffic_direction":"EGRESS","trace_observation_point":"TO_STACK","is_reply":false,"Summary":"TCP Flags: ACK"}`)
	CheckError(e)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}