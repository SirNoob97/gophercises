package main

import (
	"fmt"
	"log"
	"regexp"

	phoneDb "github.com/SirNoob97/gophercises/phone/db"
)

const (
	host     = "192.168.0.21"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbName   = "phone"
)

func main() {

	connConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s", host, port, user, password)

	must(phoneDb.Reset("pgx", connConfig, dbName))

	connConfig = fmt.Sprintf("%s dbname=%s", connConfig, dbName)
	must(phoneDb.Migrate("pgx", connConfig))

	conn, err := phoneDb.Open("pgx", connConfig)
	must(err)
	defer conn.Close()

	must(conn.Seed())

	phones, err := conn.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or Removing...", number)
			existing, err := conn.FindPhone(number)
			must(err)
			if existing != nil {
				must(conn.DeletePhone(p.ID))
			} else {
				p.Number = number
				must(conn.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes requiered")
		}
	}
}

func normalize(phone string) string {
	r := regexp.MustCompile("\\D")
	return r.ReplaceAllString(phone, "")
}

func must(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
