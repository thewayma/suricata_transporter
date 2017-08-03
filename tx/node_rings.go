package tx

import (
	"github.com/thewayma/suricata_transporter/g"
	"github.com/toolkits/consistent/rings"
    "sort"
)

func KeysOfMap(m map[string]string) []string {
    keys := make(sort.StringSlice, len(m))
    i := 0
    for key := range m {
        keys[i] = key
        i++
    }

    keys.Sort()
    return []string(keys)
}

//!< 建立一致性哈希环
func initNodeRings() {
	cfg := g.Config()

	JudgeNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Judge.Replicas), KeysOfMap(cfg.Judge.Cluster))
	//GraphNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Graph.Replicas), KeysOfMap(cfg.Graph.Cluster))
}
