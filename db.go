package main

import (
	"database/sql"
	"log"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

func setup_db() *sql.DB {
	host, password := os.Getenv("DB_HOST"), os.Getenv("DB_PASSWORD")
	serviceURI := "postgres://avnadmin:" + password + "@" + host + ":28073/defaultdb?sslmode=require"

	conn, _ := url.Parse(serviceURI)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"
	db, err := sql.Open("postgres", conn.String())

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createTable(db *sql.DB, path string) error {
	query := readSQLFile(path)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating test_table: %v", err)
		return err
	}
	log.Println("created table")
	return nil
}

func createChannelRow(db *sql.DB, channel_id string, tag string, name string) error {
	insertQuery := readSQLFile("./sql/create_channel.sql")
	_, err := db.Exec(insertQuery, channel_id, tag, name, "chess")
	if err != nil {
		return err
	}
	return nil
}

func deleteRow(db *sql.DB, tag string) error {
	query := readSQLFile("./sql/delete_row.sql")
	_, err := db.Exec(query, tag)
	if err != nil {
		return err
	}
	return nil
}

func deleteTable(db *sql.DB, path string) error {
	query := readSQLFile(path)
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("deleted table")
	return nil
}
