package main

import (
	"database/sql"
	"fmt"
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
	// Insert a row
	insertQuery := readSQLFile("./sql/create_channel.sql")
	_, err := db.Exec(insertQuery, channel_id, tag, name, "chess")
	if err != nil {
		return err
	}
	return nil
}

type Channel struct {
	id         int
	channel_id string
	tag        string
	name       string
	category   string
}

func readChannels(db *sql.DB) []Channel {
	channels := []Channel{}
	query := readSQLFile("./sql/read_all_channels.sql")
	rows, err := db.Query(query)
	if err != nil {
		err_fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var channel Channel
		err = rows.Scan(&channel.id, &channel.channel_id, &channel.tag, &channel.name, &channel.category)
		if err != nil {
			err_fatal(err)
		}
		channels = append(channels, channel)
	}
	return channels
}

func deleteRow(db *sql.DB, tag string) error {
	query := readSQLFile("./sql/delete_row.sql")
	_, err := db.Exec(query, tag)
	if err != nil {
		return err
	}
	return nil
}
func make_trigger(db *sql.DB) error {
	query := readSQLFile("./sql/trigger.sql")
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("created trigger")
	return nil
}
func deleteTable(db *sql.DB, path string) error {
	query := readSQLFile(path)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("error deleting table: %v", err)
		return err
	}
	log.Println("deleted table")
	return nil
}
