package model

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

type test struct {
	fail   bool
	data   string
	result RulesWrapper
}

var data = []test{
	{true, "", RulesWrapper{Groups{Group{Rules: nil}}}},
	{false, `{"group":{"condition":"AND","rules":[{"key":"team label","operator":"Equal","value":"my-team"},{"key":"service label","operator":"Equal","value":"thing"}]}}`,
		RulesWrapper{
			Groups{
				Group{
					Condition: "AND",
					Rules: []Rule{
						{
							Operator: "Equal",
							Key:      "team label",
							ValuesWrapper: ValuesWrapper{
								Value: []string{"my-team"},
							},
						},
						{
							Operator: "Equal",
							Key:      "service label",
							ValuesWrapper: ValuesWrapper{
								Value: []string{"thing"},
							},
						},
					},
				},
			},
		},
	},
	{false, `{"group":{"condition":"OR","rules":[]}}`, RulesWrapper{
		Groups{
			Group{
				Condition: "OR",
				Rules:     []Rule{},
			},
		},
	},
	},

	{false, `{"group":{"condition":"AND","rules":[{"key":"team label","operator":"Equal","value":"stuff"},{"key":"service label","operator":"Equal","value":"thing"}]}}`,
		RulesWrapper{
			Groups{
				Group{
					Condition: "AND",
					Rules: []Rule{
						{
							Operator: "Equal",
							Key:      "team label",
							ValuesWrapper: ValuesWrapper{
								Value: []string{"stuff"},
							},
						},
						{
							Operator: "Equal",
							Key:      "service label",
							ValuesWrapper: ValuesWrapper{
								Value: []string{"thing"},
							},
						},
					},
				},
			},
		},
	},

	{false, `{"key":"AWS Service","operator":"In","value":["AWS Support Developer"," AWS Premium Support"," AWS Support Business","AWS Config","AWS Glue"]}`,
		RulesWrapper{
			Groups{
				Group{
					Rules: []Rule{
						{
							Operator: "In",
							Key:      "AWS Service",
							ValuesWrapper: ValuesWrapper{
								Value: []string{"AWS Support Developer", " AWS Premium Support", " AWS Support Business", "AWS Config", "AWS Glue"},
							},
						},
					},
				},
			},
		},
	},

	{false, `{"key":"team label","operator":"Equal","value":"datadevils"}`,
		RulesWrapper{
			Groups{
				Group{
					Rules: []Rule{
						{
							Operator: "Equal",
							Key:      "team label",
							ValuesWrapper: ValuesWrapper{
								Value: []string{"datadevils"},
							},
						},
					},
				},
			},
		},
	},
}

func TestBusinessContext(t *testing.T) {
	assert := assert.New(t)

	for index, item := range data {
		var wrapper RulesWrapper
		err := json.Unmarshal([]byte(item.data), &wrapper)
		if item.fail {
			assert.NotNil(err, "index %d", index)
		} else {
			assert.Nil(err, "index %d", index)
		}

		assert.Equal(item.result, wrapper, "index %d", index)
	}
}

type test_data struct {
	fail   bool
	data   string
}


var data2 = []test{
	{true, `
	{
		"business_context": {
		  "name": "Business Context-1",
		  "priority": 0,
		  "criteria": {
			"condition": "OR",
			"rules": [
			  {
				"key": "Owner",
				"operator": "Equal",
				"value": "Jon Doe"
			  },
			  {
				"key": "owner_aws_tag",
				"operator": "Like",
				"value": "Jon",
				"aws_tag": true
			  },
			  {
				"group": {
				  "condition": "AND",
				  "rules": [
					{
					  "key": "AWS Account",
					  "operator": "Not Equal",
					  "value": "aws_account_id_1"
					},
					{
					  "key": "Application",
					  "operator": "Is Untagged"
					},
					{
					  "key": "AWS Service",
					  "operator": "IN",
					  "value": [
						"Amazon DynamoDB",
						"Amazon Elastic Compute Cloud",
						"Amazon RDS Service",
						"Amazon Redshift",
						"Amazon Simple Storage Service",
						"Amazon Virtual Private Cloud",
						"Amazon ElastiCache"
					  ]
					}
				  ]
				}
			  }
			]
		  }
		}
	  }
	  `}

func TestBusinessContext(t *testing.T) {
	assert := assert.New(t)

	for index, item := range data2 {
		var data BusinessContext
		err := json.Unmarshal([]byte(item.data), &data)

		fmt.Printf("%v\n", data)

	}
}
