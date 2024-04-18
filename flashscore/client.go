package flashscore

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Client struct {
	http    *http.Client
	baseUri string
	apiKey  string
}

func NewClient(baseUri, apiKey string) *Client {
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

	return &Client{c, baseUri, apiKey}
}

func (c Client) GetSchedule() (Response, error) {
	url := c.baseUri + "/v1/events/list?locale=en_GB&timezone=-4&sport_id=12&indent_days=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return Response{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "flashscore.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", c.apiKey)

	res, err := c.http.Do(req)
	if res.StatusCode == http.StatusForbidden {
		log.Println("Forbidden - wrong API key?")
		return Response{}, err
	}
	if err != nil {
		log.Println("Error executing GET request", err)
		return Response{}, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("request failed with status code: %v, req: %v\n", res.StatusCode, req)
		return Response{}, fmt.Errorf("request failed with status code: %v, body: %v", res.StatusCode, res.Body)
	}

	if res.Body == nil {
		return Response{}, errors.New("no response body")
	}
	defer res.Body.Close()

	return NewResponse(res.Body)
}
