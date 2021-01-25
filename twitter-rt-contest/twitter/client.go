package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// Client twitter client
type Client struct {
	client *http.Client
}

// New create a new twitter client
func New(consumer, secret string) (*Client, error) {
	client, err := twitterClient(consumer, secret)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func twitterClient(consumer, secret string) (*http.Client, error) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.SetBasicAuth(consumer, secret)

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		return nil, err
	}

	var conf oauth2.Config
	return conf.Client(context.Background(), &token), nil
}

// Retweeters get all the users who retweeted the tweet
func (c *Client) Retweeters(twID string) ([]string, error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", twID)
	res, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var rts []struct {
		User struct {
			ScreenName string `json:"screen_name"`
		} `json:"user"`
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&rts)
	if err != nil {
		panic(err)
	}

	retweeters := make([]string, 0, len(rts))
	for _, user := range rts {
		retweeters = append(retweeters, user.User.ScreenName)
	}

	return retweeters, nil
}
