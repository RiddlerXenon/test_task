package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	Add(key, value string, ttl time.Duration) error
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type requestPayload struct {
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

func (c *HTTPClient) Add(key, value string, ttl time.Duration) error {
	return c.request("add", key, value, ttl)
}

func (c *HTTPClient) Set(key, value string, ttl time.Duration) error {
	return c.request("set", key, value, ttl)
}

func (c *HTTPClient) Get(key string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/get?key=%s", c.baseURL, url.QueryEscape(key)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("key not found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *HTTPClient) Del(key string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/del?key=%s", c.baseURL, url.QueryEscape(key)), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("key not found")
	}
	return nil
}

func (c *HTTPClient) request(endpoint, key, value string, ttl time.Duration) error {
	url := fmt.Sprintf("%s/%s?key=%s", c.baseURL, endpoint, url.QueryEscape(key))

	payload := requestPayload{
		Value: value,
		TTL:   ttl,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to execute request")
	}
	return nil
}
