package db

import (
	"database/sql"

	// postgres driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB db struct
type DB struct {
	db *sql.DB
}

// Phone phone model
type Phone struct {
	ID     int
	Number string
}

// Open create a database connection
func Open(driverName, dataSource string) (DB, error) {
	conn, err := sql.Open(driverName, dataSource)
	if err != nil {
		return DB{}, err
	}

	return DB{conn}, nil
}

// Close close the database connection
func (db *DB) Close() error {
	return db.db.Close()
}

// Seed insert fixed data to the database
func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"(123)456-7892",
	}

	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}

	return nil
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

// AllPhones get all the records
func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT * FROM phone_number")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var ret []Phone

	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

// FindPhone find phone using the number
func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	statement := `SELECT * FROM phone_number WHERE value = $1`
	row := db.db.QueryRow(statement, number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// UpdatePhone update the phone number
func (db *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_number SET  value=$2 WHERE id=$1`
	_, err := db.db.Exec(statement, p.ID, p.Number)
	return err
}

// DeletePhone delete a record using the id
func (db *DB) DeletePhone(id int) error {
	statement := `DELETE FROM phone_number WHERE id=$1`
	_, err := db.db.Exec(statement, id)
	return err
}

// Migrate create phone_number table
func Migrate(driverName, dataSource string) error {
	conn, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	err = createPhoneNumberTable(conn)
	if err != nil {
		return err
	}

	return conn.Close()
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

// Reset reset the database status
func Reset(driverName, dataSource, dbName string) error {
	conn, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	err = resetDB(conn, dbName)
	if err != nil {
		return err
	}

	return conn.Close()
}

func resetDB(db *sql.DB, name string) error {
	statement := "DROP DATABASE IF EXISTS " + name
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	statement := "CREATE DATABASE " + name
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

// get phone number by id
//func GetPhone(db *sql.DB, id int) (string, error) {
//var number string
//statement := `SELECT * FROM phone_number WHERE id = $1`
//row := db.QueryRow(statement, id)
//err := row.Scan(&id, &number)
//if err != nil {
//return "", err
//}
//return number, nil
//}
