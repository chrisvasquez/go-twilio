package twilio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	DEFAULT_TIMEOUT = 30000
	API_BASE        = "https://api.twilio.com/"
	API_VERSION     = "2010-04-01"
)

type client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	httpTimeout time.Duration
}

// NewClient return a new Twilio HTTP client
func NewClient(apiKey, apiSecret string) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, DEFAULT_TIMEOUT * time.Millisecond}
}

// NewClientWithCustomHttpConfig returns a new Twilio HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(apiKey, apiSecret string, httpClient *http.Client) (c *client) {
	timeout := httpClient.Timeout
	if timeout <= 0 {
		timeout = DEFAULT_TIMEOUT * time.Millisecond
	}
	return &client{apiKey, apiSecret, httpClient, timeout}
}

// NewClient returns a new Twilio HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, timeout}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on Twilio API")
	}
}

// do prepare and process HTTP request to Twilio API
func (c *client) do(payload string) (response []byte, err error) {

	connectTimer := time.NewTimer(c.httpTimeout)
	var urlStr = fmt.Sprintf("%s%s/Accounts/%s/Messages.json", API_BASE, API_VERSION, c.apiKey)

	rb := *strings.NewReader(payload)

	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		return
	}
	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}

	return response, err

}
