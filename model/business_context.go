package model

import (
	"encoding/json"
)

// BusinessContextsResult results for multiple contexts
type BusinessContextsResult struct {
	Status          string            `json:"status,omitempty"`
	BusinessContext []BusinessContext `json:"business_context,omitempty"`
}

// CreateBusinessContext input to the create call
type CreateBusinessContext struct {
	BusinessContext InputBusinessContext `json:"business_context,omitempty"`
}

// BusinessContextResult is the result from the list call
type BusinessContextResult struct {
	Status          string          `json:"status,omitempty"`
	BusinessContext BusinessContext `json:"business_context,omitempty"`
}

// UpdateBusinessContextResult is the result from the update call
type UpdateBusinessContextResult struct {
	BusinessContext BusinessContext `json:"business_context,omitempty"`
}

// BusinessContext struct for BusinessContext
type BusinessContext struct {
	UUID       string   `json:"uuid,omitempty"`
	Name       string   `json:"name,omitempty"`
	ParentUUID string   `json:"parent_uuid,omitempty"`
	ParentPath []string `json:"parent_path,omitempty"`
	Priority   int32    `json:"priority,omitempty"`
	Criteria   Criteria `json:"criteria,omitempty"`
}

// Criteria struct for Criteria
type Criteria struct {
	Rules     []RulesWrapper `json:"rules,omitempty"`
	Condition string         `json:"condition,omitempty"`
}

// RulesWrapper to modify rule unmarshalling
type RulesWrapper struct {
	Groups
}

// Groups struct
type Groups struct {
	Group Group `json:"group,omitempty"`
}

// Group struct
type Group struct {
	Condition string `json:"condition,omitempty"`
	Rules     []Rule `json:"rules,omitempty"`
}

// Rule struct
type Rule struct {
	Operator      string        `json:"operator,omitempty"`
	Key           string        `json:"key,omitempty"`
	ValuesWrapper ValuesWrapper `json:"value,omitempty"`
}

// ValuesWrapper to modify Values unmarshalling
type ValuesWrapper struct {
	Value []string
}

// UnmarshalJSON for modified format
func (w *RulesWrapper) UnmarshalJSON(data []byte) error {
	var rule Rule

	err := json.Unmarshal(data, &rule)
	w.Group.Rules = []Rule{rule}
	if err != nil || rule.Key == "" {
		err := json.Unmarshal(data, &w.Groups)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalJSON for modified format
func (w *ValuesWrapper) UnmarshalJSON(data []byte) error {

	err := json.Unmarshal(data, &w.Value)
	if err != nil {
		var str string
		err := json.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		w.Value = append(w.Value, str)
	}

	return nil
}

// InputBusinessContext for input
type InputBusinessContext struct {
	UUID       string        `json:"uuid,omitempty"`
	Name       string        `json:"name,omitempty"`
	ParentUUID string        `json:"parent_uuid,omitempty"`
	ParentPath []string      `json:"parent_path,omitempty"`
	Priority   int32         `json:"priority,omitempty"`
	Criteria   InputCriteria `json:"criteria,omitempty"`
}

// Criteria struct for Criteria
type InputCriteria struct {
	Rules     []InputGroups `json:"rules,omitempty"`
	Condition string        `json:"condition,omitempty"`
}

type InputGroups struct {
	Group InputGroup `json:"group,omitempty"`
}

type InputGroup struct {
	Condition string     `json:"condition,omitempty"`
	Rules     []ShitRule `json:"rules,omitempty"`
}

type ShitRule struct {
	Operator string   `json:"operator,omitempty"`
	Key      string   `json:"key,omitempty"`
	Values   []string `json:"value,omitempty"`
}
