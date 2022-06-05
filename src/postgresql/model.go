package postgresql

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Psqlconn struct {
	Db *sql.DB
	Log *logrus.Logger
}
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
