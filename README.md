# What is Hubble-Storage?

The cilium traffics is very large, so if you always store each traffic, the database will exceed the allocated resources.

Recognizing this problem, Hubble-Storage will process the json traffic logs consumed from Kafka, calculate the corresponding hash code, normalize and store it in the postgres database. Based on the previously saved hash code, Hubble-Storage will only update the corresponding quantity in the database.

`Hubble-Storage` and [Hubble-Log-Collector](https://github.com/CyberSecurityBox/hubble-log-collector) are used together.

For example:
* Traffic log saved:

```
id | verdict | ip_source | ... | flow_count | checksum | cluster_name |
1  | DROPPED | 10.10.1.2 | ... |     1      |   ABCD   | test-cilium  |
```

* There are 10 same traffic logs:

```
id | verdict | ip_source | ... | flow_count | checksum | cluster_name |
1  | DROPPED | 10.10.1.2 | ... |     10     |   ABCD   | test-cilium  |
```

# How does it work?
## Steps:
1. Consume messages from Kafka.
2. Parse json message to `Flow` struct.
```
type Flow struct {
	Verdict              string `json:"verdict"`
	IpSource             string `json:"ip_source"`
	IpDestination        string `json:"ip_destination"`
	L4Protocol           string `json:"l_4_protocol"`
	L4SourcePort         int    `json:"l_4_source_port"`
	L4DestinationPort    int    `json:"l_4_destination_port"`
	SourcePodName        string `json:"source_pod_name"`
	SourceNamespace      string `json:"source_namespace"`
	DestinationPodName   string `json:"destination_pod_name"`
	DestinationNamespace string `json:"destination_namespace"`
	Layer                string `json:"layer"`
	Checksum             string `json:"checksum"`
	FlowCount            int    `json:"flow_count"`
	Metadata             string `json:"metadata"`
	ClusterName          string `json:"cluster_name"`
}
```
3. Calculate the hash code based on the parameters that change the least.
```
var flow *Flow

flow.Checksum = utils.FlowHash([]string{
		flow.Verdict,
		flow.IpSource,
		flow.IpDestination,
		flow.L4Protocol,
		fmt.Sprintf("%d", flow.L4DestinationPort),
		flow.SourcePodName,
		flow.SourceNamespace,
		flow.DestinationPodName,
		flow.DestinationNamespace,
		flow.Layer,
		flow.ClusterName,
	})
```
4. Based on the hash code, should insert a new record to the postgres database or update the flow_count of the corresponding record.
```
err = psql.UpdateFlowCount(flow.Checksum)
if err != nil {
	err = psql.Insert(*flow)
	_ = slack.SendWithVerdict(flow, []string{"DROPPED", "ERROR"}, core.Log)
	if err != nil {
		return err
	}
}
```
5. Push notification on the Slack channel if the verdict is `DROPPED` or `ERROR`.
```
_ = slack.SendWithVerdict(flow, []string{"DROPPED", "ERROR"}, core.Log)
```

# How to use it?
1. Let's create a namespace named hubble-dev
```
kubectl create namespace hubble-dev
```
2. Create a secret file containing connection parameters to Kafka, Postgresql and Slack webhook.
```
cat <<EOF > hubble-storage.yaml
apiVersion: v1
data:
  KAFKA_BROKERS: <base64>
  KAFKA_CONSUME_GROUP: <base64>
  KAFKA_PASS: <base64>
  KAFKA_TOPIC: <base64>
  KAFKA_USER: <base64>
  POSTGRES_DATABASE: <base64>
  POSTGRES_PASS: <base64>
  POSTGRES_PORT: <base64>
  POSTGRES_SERVER: <base64>
  POSTGRES_USER: <base64>
  SLACK_WEBHOOK_URL: <base64>
kind: Secret
metadata:
  name: hubble-storage
  namespace: hubble-dev
type: Opaque
EOF
```
2. Apply [secret](.k8s/secret.yaml) to cluster
```
kubectl apply -f hubble-storage.yaml
```
3. Create a [Deployment](.k8s/deployment.yaml) to perform traffic log processing.
```
kubectl apply -f .k8s/deployment.yaml
```
## Notes
* You should only set the variable value of `KAFKA_CONSUME_GROUP` as unique, it should not be changed, because it may cause the workload to consume all the old traffic log again.
* Create a Postgres database with a structure like [flows.sql](.postgresql/flows.sql)