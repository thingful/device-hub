// Copyright © 2017 thingful

package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	hub "github.com/thingful/device-hub"
)

type httpout struct {
	url    string
	client *http.Client
}

func NewHTTPEndpoint(url string, timeOutInMS int) httpout {

	return httpout{
		url: url,
		client: &http.Client{
			Timeout: time.Millisecond * time.Duration(timeOutInMS),
		},
	}

}

func (h httpout) Write(message hub.Message) error {

	j, err := json.Marshal(message)

	if err != nil {
		return err
	}

	resp, err := h.client.Post(h.url, "application/json", bytes.NewBuffer(j))

	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("unexpected response %s", resp.StatusCode)

}
