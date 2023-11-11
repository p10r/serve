package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Client struct {
	http    *http.Client
	fullUrl string
}

func NewClient(fullUrl string) *Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return &Client{c, fullUrl}
}

func (c Client) SendMessage(msg Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	res, err := http.Post(c.fullUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		log.Printf("Discord request failed with status code: %v\n", res.StatusCode)
		return fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	return nil
}
