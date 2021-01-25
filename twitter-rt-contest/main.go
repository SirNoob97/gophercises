package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var (
		keysFile  string
		usersFile string
		tweetID   string
	)

	flag.StringVar(&keysFile, "keys", ".keys.json", "The file where you store your consumer key and secret for the twitter API.")
	flag.StringVar(&usersFile, "file", "users.csv", "The file where users who have retweeted the tweet are stored. This will be created if does not exist.")
	flag.StringVar(&tweetID, "tweet", "1353417790804385792", "The ID of the tweet you wish to find retweeters of.")
	flag.Parse()

	consumer, secret, err := keys(keysFile)
	if err != nil {
		log.Fatalf("Could not find the file %s", keysFile)
	}

	var client http.Client
	accessToken, err := twitterToken(&client, consumer, secret)
	if err != nil {
		log.Fatal(err)
	}
	newUsersnames, err := retweeters(&client, accessToken, tweetID)
	if err != nil {
		log.Fatal(err)
	}

	existUsernames := existing(usersFile)
	allUsernames := merge(newUsersnames, existUsernames)

	err = persistUser(usersFile, allUsernames)
	if err != nil {
		log.Fatal(err)
	}
}

func keys(keysFile string) (string, string, error) {
	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}

	file, err := os.Open(keysFile)
	defer file.Close()
	if err != nil {
		return "", "", err
	}

	dec := json.NewDecoder(file)
	dec.Decode(&keys)

	return keys.Key, keys.Secret, nil
}

func twitterToken(client *http.Client, consumer, secret string) (string, error) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.SetBasicAuth(consumer, secret)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var token struct {
		Type  string `json:"token_type"`
		Token string `json:"access_token"`
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		return "", err
	}
	accessToken := "Bearer " + token.Token
	return accessToken, nil
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

func existing(usersFile string) []string {
	f, err := os.Open(usersFile)
	if err != nil {
		return []string{}
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	users := make([]string, 0, len(lines))

	for _, line := range lines {
		users = append(users, line[0])
	}
	return users
}

func merge(a, b []string) []string {
	uniq := make(map[string]struct{}, 0)
	for _, user := range a {
		uniq[user] = struct{}{}
	}
	for _, user := range b {
		uniq[user] = struct{}{}
	}
	ret := make([]string, 0, len(uniq))
	for user := range uniq {
		ret = append(ret, user)
	}

	return ret
}

func persistUser(usersFile string, users []string) error {
	f, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	for _, user := range users {
		if err := w.Write([]string{user}); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}
