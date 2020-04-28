package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sstarcher/yotascale-sdk-golang/model"
)

// Client to access okta
type Client struct {
	client *http.Client
	url    string
	token  string
}

// NewClient creates a yotascale api client
func NewClient() (*Client, error) {
	token, exists := os.LookupEnv("YOTASCALE_TOKEN")
	if !exists {
		return nil, errors.New("YOTASCALE_TOKEN environment variable must be set")
	}

	url := "https://app-api.yotascale.io/api/v1/"
	client := &Client{
		client: &http.Client{},
		url:    url,
		token:  token,
	}

	return client, nil
}

// ListContexts for all business contexts
func (c *Client) ListContexts() ([]model.BusinessContext, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"business-contexts", nil)
	req.Header.Add("x-api-key", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var body []byte
	if resp.StatusCode != 200 {
		debugResponse(req, resp)

		body, _ = ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	body, _ = ioutil.ReadAll(resp.Body)
	var data model.BusinessContextsResult
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data.BusinessContext, nil
}

// CreateContext for a new business context
func (c *Client) CreateContext(parentUUID string, item model.InputBusinessContext) (*model.BusinessContext, error) {
	data := model.CreateBusinessContext{item}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url+"business-contexts", bytes.NewBuffer((dataBytes)))
	req.Header.Add("x-api-key", c.token)
	req.Header.Add("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("parent_uuid", parentUUID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var body []byte
	if resp.StatusCode != 200 {
		debugResponse(req, resp)
		body, _ = ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	body, _ = ioutil.ReadAll(resp.Body)
	var result model.BusinessContextResult
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &result.BusinessContext, nil
}

// UpdateContext for the existing business context
func (c *Client) UpdateContext(item model.InputBusinessContext) (*model.BusinessContext, error) {
	data := model.CreateBusinessContext{item}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.url+"business-contexts/"+item.UUID, bytes.NewBuffer(jsonData))
	req.Header.Add("x-api-key", c.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var body []byte
	if resp.StatusCode != 200 {
		debugResponse(req, resp)
		body, _ = ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	body, _ = ioutil.ReadAll(resp.Body)
	var result model.UpdateBusinessContextResult
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &result.BusinessContext, nil
}

// DeleteContext removes the given business context
func (c *Client) DeleteContext(uuid string) error {
	req, err := http.NewRequest(http.MethodDelete, c.url+"business-contexts/"+uuid, nil)
	req.Header.Add("x-api-key", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)

		var rsp error
		if err = json.Unmarshal(body, &rsp); err != nil {
			debugResponse(req, resp)
			return errors.New(string(body))
		}
		return rsp
	}

	return nil
}

func debugResponse(req *http.Request, resp *http.Response) {
	if log.GetLevel() >= log.DebugLevel {
		dump, err := httputil.DumpRequest(req, true)
		log.Debug("**************************************************")
		log.Debugf("Request %s\n", string(dump))
		if err != nil {
			log.Debugf("\terr %v\n", err)
		}

		log.Debug("#################################################")

		dump, err = httputil.DumpResponse(resp, true)
		log.Debugf("\n\nResponse%s\n", dump)
		if err != nil {
			log.Debugf("\terr %v\n", err)
		}
		log.Debug("**************************************************")
	}
}
