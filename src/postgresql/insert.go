package postgresql

func (psqlconn *Psqlconn) Insert(flow Flow) error {
	var insertDynStmt = "insert into " +
		"flows(verdict, ip_source, ip_destination, l4_protocol, l4_source_port, l4_destination_port, source_pod_name, source_namespace, destination_pod_name, destination_namespace, layer, checksum, flow_count, metadata, cluster_name) " +
		"values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)"
	var err error
	var in = []interface{}{
		flow.Verdict,
		flow.IpSource,
		flow.IpDestination,
		flow.L4Protocol,
		flow.L4SourcePort,
		flow.L4DestinationPort,
		flow.SourcePodName,
		flow.SourceNamespace,
		flow.DestinationPodName,
		flow.DestinationNamespace,
		flow.Layer,
		flow.Checksum,
		flow.FlowCount,
		flow.Metadata,
		flow.ClusterName,
	}
	_, err = psqlconn.Db.Exec(insertDynStmt, in...)
	if err != nil {
		psqlconn.Log.WithError(err).Error("Can't insert db.Exec")
		return err
	}

	return nil
}
