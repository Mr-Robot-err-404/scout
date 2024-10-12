package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Quota struct {
	id        int
	quota     int
	timestamp time.Time
}
type OAuth struct {
	Access_token string
}

func refresh_quota(db *sql.DB) {
	refresh_token := os.Getenv("REFRESH_TOKEN")
	client_id, client_secret := os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")

	now := time.Now().Unix()
	last_update := read_quota(db).timestamp.Unix()
	diff := now - last_update

	if diff < 30 {
		return
	}
	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("refresh_token", refresh_token)
	data.Set("grant_type", "refresh_token")

	route := "https://www.googleapis.com/oauth2/v4/token"
	resp, err := http.PostForm(route, data)

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	var s OAuth
	err = json.Unmarshal(body, &s)

	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("ACCESS_TOKEN", s.Access_token)

	if err != nil {
		log.Fatal(err)
	}
	query := readSQLFile("./sql/refresh.sql")
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("refreshed access token")
	fmt.Println(s.Access_token)
}

func read_quota(db *sql.DB) Quota {
	var quota Quota
	query := readSQLFile("./sql/read_quota.sql")
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&quota.id, &quota.quota, &quota.timestamp)
		if err != nil {
			log.Fatal(err)
		}
	}
	return quota
}
