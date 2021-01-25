package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/SirNoob97/gophercises/twitter-rt-contest/twitter"
)

func main() {
	var (
		keysFile  string
		usersFile string
		tweetID   string
		nWinners  int
	)

	flag.StringVar(&keysFile, "keys", ".keys.json", "The file where you store your consumer key and secret for the twitter API.")
	flag.StringVar(&usersFile, "file", "users.csv", "The file where users who have retweeted the tweet are stored. This will be created if does not exist.")
	flag.StringVar(&tweetID, "tweet", "1353417790804385792", "The ID of the tweet you wish to find retweeters of.")
	flag.IntVar(&nWinners, "winners", 0, "The number of winners to pick for the contest.")
	flag.Parse()

	consumer, secret, err := keys(keysFile)
	if err != nil {
		log.Fatalf("Could not find the file %s", keysFile)
	}

	client, err := twitter.New(consumer, secret)
	if err != nil {
		log.Fatal(err)
	}
	newUsersnames, err := client.Retweeters(tweetID)
	if err != nil {
		log.Fatal(err)
	}

	existUsernames := existing(usersFile)
	allUsernames := merge(newUsersnames, existUsernames)

	err = persistUser(usersFile, allUsernames)
	if err != nil {
		log.Fatal(err)
	}

	if nWinners != 0 {
		existUsernames = existing(usersFile)
		winners := pickWinners(existUsernames, nWinners)
		for _, user := range winners {
			fmt.Printf("  * %s\n", user)
		}
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

func pickWinners(users []string, nWinners int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(users))
	winners := perm[:nWinners]
	ret := make([]string, 0, nWinners)
	for _, idx := range winners {
		ret = append(ret, users[idx])
	}
	return ret
}
