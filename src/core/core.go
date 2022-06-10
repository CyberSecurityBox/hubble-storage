package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"hubble-storage/src/postgresql"
	"hubble-storage/src/slack"
	"hubble-storage/src/utils"
)

func (core *Core) parser(key string, value string) (*postgresql.Flow, error) {
	var flow postgresql.Flow
	var parse map[string]interface{}
	var err error

	flow.ClusterName = key

	err = json.Unmarshal([]byte(value), &parse)
	if err != nil {
		core.Log.WithError(err).Error("Can't parse")
		return &flow, err
	}

	if parse["verdict"] != nil {
		flow.Verdict = parse["verdict"].(string)
	} else {
		core.Log.Info(value)
		flow.Verdict = ""
	}

	if parse["IP"].(map[string]interface{})["source"] != nil {
		flow.IpSource = parse["IP"].(map[string]interface{})["source"].(string)
	} else {
		core.Log.Info(value)
		flow.IpSource = ""
	}

	if parse["IP"].(map[string]interface{})["destination"] != nil {
		flow.IpDestination = parse["IP"].(map[string]interface{})["destination"].(string)
	} else {
		core.Log.Info(value)
		flow.IpDestination = ""
	}

	if parse["l4"] != nil {
		for k, v := range parse["l4"].(map[string]interface{}) {
			flow.L4Protocol = k
			if v.(map[string]interface{})["source_port"] != nil {
				flow.L4SourcePort = int(v.(map[string]interface{})["source_port"].(float64))
			} else {
				core.Log.Info(value)
				flow.L4SourcePort = 0
			}
			if v.(map[string]interface{})["destination_port"] != nil {
				flow.L4DestinationPort = int(v.(map[string]interface{})["destination_port"].(float64))
			} else {
				core.Log.Info(value)
				flow.L4DestinationPort = 0
			}
		}
	} else {
		core.Log.Info(value)
		flow.L4Protocol = ""
		flow.L4SourcePort = 0
		flow.L4DestinationPort = 0
	}

	if parse["source"] != nil {
		source := parse["source"].(map[string]interface{})
		identity := int(source["identity"].(float64))
		labels := source["labels"].([]interface{})
		if identity == 0 || identity == 1 || identity == 2 || identity == 3 ||
			identity == 4 || identity == 5 || identity == 6 {
			flow.SourcePodName = labels[0].(string)
			flow.SourceNamespace = ""
		} else {
			if source["pod_name"] != nil {
				flow.SourcePodName = source["pod_name"].(string)
				flow.SourceNamespace = source["namespace"].(string)
			} else {
				if i:=0; i<len(labels) {
					flow.SourcePodName = labels[0].(string)
					flow.SourceNamespace = ""
				} else {
					core.Log.Info(value)
					flow.SourcePodName = ""
					flow.SourceNamespace = ""
				}
			}
		}
	} else {
		core.Log.Info(value)
		flow.SourcePodName = ""
		flow.SourceNamespace = ""
	}

	if parse["destination"] != nil {
		destination := parse["destination"].(map[string]interface{})
		identity := int(destination["identity"].(float64))
		labels := destination["labels"].([]interface{})
		if identity == 0 || identity == 1 || identity == 2 || identity == 3 ||
			identity == 4 || identity == 5 || identity == 6 {
			flow.DestinationPodName = labels[0].(string)
			flow.DestinationNamespace = ""
		} else {
			if destination["pod_name"] != nil {
				flow.DestinationPodName = destination["pod_name"].(string)
				flow.DestinationNamespace = destination["namespace"].(string)
			} else {
				if i:=0; i<len(labels) {
					flow.DestinationPodName = labels[0].(string)
					flow.DestinationNamespace = ""
				} else {
					core.Log.Info(value)
					flow.DestinationPodName = ""
					flow.DestinationNamespace = ""
				}
			}
		}
	} else {
		core.Log.Info(value)
		flow.DestinationPodName = ""
		flow.DestinationNamespace = ""
	}

	if parse["Type"] != nil {
		flow.Layer = parse["Type"].(string)
	} else {
		core.Log.Info(value)
		flow.Layer = ""
	}

	flow.Metadata = value

	core.Log.Info(fmt.Sprintf("%s: %s", flow.Verdict, flow.SourcePodName))

	return &flow, nil
}
func (core *Core) store(flow *postgresql.Flow) error {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		core.Host,
		core.Port,
		core.User,
		core.Password,
		core.Dbname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		core.Log.WithError(err).Error("Can't insert db.Exec")
		return err
	}

	defer db.Close()

	psql := postgresql.Psqlconn{
		Db:  db,
		Log: core.Log,
	}

	flow.Checksum = utils.FlowHash([]string{
		flow.Verdict,
		flow.IpSource,
		flow.IpDestination,
		flow.L4Protocol,
		//fmt.Sprintf("%d", flow.L4SourcePort),
		fmt.Sprintf("%d", flow.L4DestinationPort),
		flow.SourcePodName,
		flow.SourceNamespace,
		flow.DestinationPodName,
		flow.DestinationNamespace,
		flow.Layer,
		flow.ClusterName,
	})

	err = psql.UpdateFlowCount(flow.Checksum)
	if err != nil {
		err = psql.Insert(*flow)
		_ = slack.SendWithVerdict(flow, []string{"DROPPED", "ERROR"}, core.Log)
		if err != nil {
			return err
		}
	}

	return nil
}
func (core *Core) MainHandler(topic string, key []byte, message []byte) error {
	var flow *postgresql.Flow
	var err error

	flow, err = core.parser(string(key), string(message))
	if err != nil {
		return err
	}

	err = core.store(flow)
	if err != nil {
		return err
	}

	return nil
}