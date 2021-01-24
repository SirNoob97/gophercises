package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type keys struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

type token struct {
	Type  string `json:"token_type"`
	Token string `json:"access_token"`
}

func main() {
	file, err := os.Open(".keys.json")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	var keys keys
	dec := json.NewDecoder(file)
	dec.Decode(&keys)

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.SetBasicAuth(keys.Key, keys.Secret)

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var token token
	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		panic(err)
	}

	accessToken := "Bearer " + token.Token
	retweeters, err := retweeters(&client, accessToken, "1353417790804385792")
	if err != nil {
		panic(err)
	}
	fmt.Println(retweeters)
}

func retweeters(client *http.Client, accessToken, twID string) ([]string, error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", twID)
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", accessToken)
	res, err := client.Do(req)
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
