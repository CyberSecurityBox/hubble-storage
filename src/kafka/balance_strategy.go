package kafka

import (
	"github.com/Shopify/sarama"
	"hubble-storage/src/env"
	"sort"
)

type myBalanceStrategy struct {
	name string
}

func (s *myBalanceStrategy) Name() string { return s.name }

func (s *myBalanceStrategy) Plan(members map[string]sarama.ConsumerGroupMemberMetadata, topics map[string][]int32) (sarama.BalanceStrategyPlan, error) {
	// Build members
	var memberIDs []string
	for memberID := range members {
		memberIDs = append(memberIDs, memberID)
	}

	sort.Sort(&balanceStrategySortable{
		memberIDs: memberIDs,
	})

	//Assemble plan
	memberPos := 0
	plan := make(sarama.BalanceStrategyPlan, len(members))
	for topic, partitions := range topics {
		for _, partition := range partitions {
			if memberPos < len(memberIDs) {
				plan.Add(memberIDs[memberPos], topic, partition)
				memberPos++
			}
		}
	}
	return plan, nil
}

func (s *myBalanceStrategy) AssignmentData(memberID string, topics map[string][]int32, generationID int32) ([]byte, error) {
	return nil, nil
}

type balanceStrategySortable struct {
	memberIDs []string
}

func (p balanceStrategySortable) Len() int { return len(p.memberIDs) }
func (p balanceStrategySortable) Swap(i, j int) {
	p.memberIDs[i], p.memberIDs[j] = p.memberIDs[j], p.memberIDs[i]
}

func (p balanceStrategySortable) Less(i, j int) bool {
	return balanceStrategyHashValue(p.memberIDs[i]) < balanceStrategyHashValue(p.memberIDs[j])
}

func balanceStrategyHashValue(vv ...string) uint32 {
	h := uint32(2166136261)
	for _, s := range vv {
		for _, c := range s {
			h ^= uint32(c)
			h *= 16777619
		}
	}
	return h
}

var MyBalanceStrategy = &myBalanceStrategy{
	name: env.KafkaGroup,
}


