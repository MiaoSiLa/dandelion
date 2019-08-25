package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tengattack/dandelion/client"
	"github.com/tengattack/dandelion/cmd/dandelion/config"
)

// Client for send events to webhook
type Client struct {
	url        string
	httpClient *http.Client
}

// Send webhook events
func (c *Client) Send(v interface{}) error {
	if c.url == "" {
		// disabled
		return nil
	}

	reqBody, err := json.Marshal(v)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client.InitHTTPRequest(req, false)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// just read all to reuse connection
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook response status %d", res.StatusCode)
	}

	return nil
}

// NewClient creates a new webhook client
func NewClient(conf *config.SectionWebhook) *Client {
	c := new(Client)
	c.url = conf.URL
	c.httpClient = &http.Client{}

	return c
}
