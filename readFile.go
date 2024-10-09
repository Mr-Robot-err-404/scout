package main

import (
	"encoding/json"
	"log"
	"os"
)

type Credentials struct {
	Client_id     string
	Client_secret string
	Refresh_token string
}

func readSQLFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func readCredentialsFile(path string) Credentials {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var credentials Credentials
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		log.Fatal(err)
	}
	return credentials
}
