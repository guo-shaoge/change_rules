package main

import (
	"fmt"
	"encoding/json"
	"os"
)

type Rules []Rule
type RuleGroup struct {
	ID       string `json:"id,omitempty"`
	Index    int    `json:"index,omitempty"`
	Override bool   `json:"override,omitempty"`
}
type LabelConstraintOp string

const (
	// In restricts the store label value should in the value list.
	// If label does not exist, `in` is always false.
	In LabelConstraintOp = "in"
	// NotIn restricts the store label value should not in the value list.
	// If label does not exist, `notIn` is always true.
	NotIn LabelConstraintOp = "notIn"
	// Exists restricts the store should have the label.
	Exists LabelConstraintOp = "exists"
	// NotExists restricts the store should not have the label.
	NotExists LabelConstraintOp = "notExists"
)

type LabelConstraint struct {
	Key    string            `json:"key,omitempty"`
	Op     LabelConstraintOp `json:"op,omitempty"`
	Values []string          `json:"values,omitempty"`
}
type PeerRoleType string

const (
	// Voter can either match a leader peer or follower peer
	Voter PeerRoleType = "voter"
	// Leader matches a leader.
	Leader PeerRoleType = "leader"
	// Follower matches a follower.
	Follower PeerRoleType = "follower"
	// Learner matches a learner.
	Learner PeerRoleType = "learner"
)

type Rule struct {
	GroupID          string            `json:"group_id"`                    // mark the source that add the rule
	ID               string            `json:"id"`                          // unique ID within a group
	Index            int               `json:"index,omitempty"`             // rule apply order in a group, rule with less ID is applied first when indexes are equal
	Override         bool              `json:"-"`          // when it is true, all rules with less indexes are disabled
	StartKey         []byte            `json:"-"`                           // range start key
	StartKeyHex      string            `json:"-"`                   // hex format start key, for marshal/unmarshal
	EndKey           []byte            `json:"-"`                           // range end key
	EndKeyHex        string            `json:"-"`                     // hex format end key, for marshal/unmarshal
	Role             PeerRoleType      `json:"-"`                        // expected role of the peers
	IsWitness        bool              `json:"-"`                  // when it is true, it means the role is also a witness
	Count            int               `json:"count"`                       // expected count of the peers
	LabelConstraints []LabelConstraint `json:"label_constraints,omitempty"` // used to select stores to place peers
	LocationLabels   []string          `json:"location_labels,omitempty"`   // used to make peers isolated physically
	IsolationLevel   string            `json:"isolation_level,omitempty"`   // used to isolate replicas explicitly and forcibly
	Version          uint64            `json:"-"`           // only set at runtime, add 1 each time rules updated, begin from 0.
	CreateTimestamp  uint64            `json:"-"`  // only set at runtime, recorded rule create timestamp
	group            *RuleGroup        // only set at runtime, no need to {,un}marshal or persist.
}

const (
	LabelKeyEngineRole = "engine_role"
	LabelValueEngineRoleWrite = "write"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("Usage: %v cur_rules.json", os.Args[0]))
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var curRules Rules
	if err := json.Unmarshal(data, &curRules); err != nil {
		panic(err)
	}

	var newRules Rules
	for _, rule := range curRules {
		var newRule Rule
		newRule.GroupID = "enable_s3_wn_region"
		newRule.ID = rule.ID
		newRules = append(newRules, newRule)
	}
	newData, err := json.MarshalIndent(newRules, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(newData))
}