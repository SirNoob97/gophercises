package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	host     = "192.168.0.21"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	db       = "phone"
)

func main() {

	connConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s", host, port, user, password)

	//conn, err := sql.Open("pgx", connConfig)
	//must(err)

	//err = resetDB(conn, db)
	//must(err)
	//conn.Close()

	connConfig = fmt.Sprintf("%s dbname=%s", connConfig, db)

  conn, err := sql.Open("pgx", connConfig)
	must(err)
	defer conn.Close()

	must(createPhoneNumberTable(conn))
  id, err := insertPhone(conn, "0123456789")
  must(err)
  fmt.Println("id=", id)
}

func createPhoneNumberTable(db *sql.DB) error {
	statement := `
    CREATE TABLE IF NOT EXISTS phone_number (
      id SERIAL,
      value VARCHAR(255)
    )`
	_, err := db.Exec(statement)
	return err
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_number(value) VALUES($1) RETURNING id`
	var id int

	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func resetDB(db *sql.DB, name string) error {
	statement := "DROP DATABASE IF EXISTS " + name
	_, err := db.Exec(statement)
	must(err)
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	statement := "CREATE DATABASE " + name
	_, err := db.Exec(statement)
	must(err)
	return nil
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

//func normalize(phone string) string{
//var cero = '0'
//var nine = '9'
//var buff bytes.Buffer

//for _, ch := range phone {
//if ch >= cero && ch <= nine {
//buff.WriteRune(ch)
//}
//}

//return buff.String()
//}
