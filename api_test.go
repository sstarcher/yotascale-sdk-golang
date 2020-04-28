package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/sstarcher/yotascale-sdk-golang/api/model"
	"github.com/stretchr/testify/assert"
)

func TestExportContexts(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient()
	assert.Nil(err)

	data, err := client.ListContexts()
	assert.Nil(err)

	var root string
	for _, context := range data {
		if context.ParentUUID == "" {
			root = context.UUID
			break
		}
	}
	process(root, data)
}

func process(id string, data []model.BusinessContext) {
	for _, item := range data {
		if item.ParentUUID == id {
			process(item.UUID, data)
			fmt.Printf("%s/%s - %s\n", strings.Join(item.ParentPath, "/"), item.Name, item.UUID)
		}
	}
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient()
	assert.Nil(err)

	data := model.InputBusinessContext{
		Name: "Unallocated",
		UUID: "1232132312312",
		Criteria: model.InputCriteria{
			Condition: "AND",
			Rules: []model.InputGroups{
				{
					Group: model.InputGroup{
						Condition: "OR",
						Rules: []model.ShitRule{
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"blah"},
							},
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"yep"},
							},
						},
					},
				},
				{
					Group: model.InputGroup{
						Condition: "OR",
						Rules: []model.ShitRule{
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"blah"},
							},
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"yep"},
							},
						},
					},
				},
			},
		},
	}

	result, err := client.UpdateContext(data)
	assert.Nil(err)
	assert.NotNil(result)
}

func TestExportTerraform(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient()
	assert.Nil(err)

	data, err := client.ListContexts()
	assert.Nil(err)

	result := map[string][]string{}
	for _, item := range data {
		if len(item.ParentPath) == 1 {
			result[item.Name] = []string{}
		}
	}

	for _, item := range data {
		if len(item.ParentPath) == 2 {
			result[item.ParentPath[1]] = append(result[item.ParentPath[1]], item.Name)
		}
	}

	myjson, err := json.Marshal(result)
	assert.Nil(err)

	for _, item := range data {
		if item.Name == "Shared" {
			fmt.Printf("terraform import 'yotascale_business_context.shared' %s\n", item.UUID)
		} else if item.Name == "Unallocated" {
			fmt.Printf("terraform import 'yotascale_business_context.unallocated[\"%s\"]' %s\n", item.ParentPath[len(item.ParentPath)-1], item.UUID)
		} else if len(item.ParentPath) == 2 {
			fmt.Printf("terraform import 'yotascale_business_context.services[\"%s.%s\"]' %s\n", item.ParentPath[len(item.ParentPath)-1], item.Name, item.UUID)
		} else if len(item.ParentPath) == 1 {
			fmt.Printf("terraform import 'yotascale_business_context.teams[\"%s\"]' %s\n", item.Name, item.UUID)
		}
	}

}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient()
	assert.Nil(err)
	err = client.DeleteContext("1232312312")
	assert.Nil(err)

}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient()
	assert.Nil(err)

	data := model.InputBusinessContext{
		Name: "create",
		Criteria: model.InputCriteria{
			Condition: "AND",
			Rules: []model.InputGroups{
				{
					Group: model.InputGroup{
						Condition: "OR",
						Rules: []model.ShitRule{
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"blah"},
							},
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"yep"},
							},
						},
					},
				},
				{
					Group: model.InputGroup{
						Condition: "OR",
						Rules: []model.ShitRule{
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"blah"},
							},
							{
								Key:      "team",
								Operator: "Equal",
								Values:   []string{"yep"},
							},
						},
					},
				},
			},
		},
	}

	result, err := client.CreateContext("12321321321312", data)
	assert.NotNil(result)
	assert.Nil(err)
}
