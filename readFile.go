package main

import (
	"encoding/json"
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
		err_fatal(err)
	}
	return string(data)
}

func read_credentials_file(path string) Credentials {
	data, err := os.ReadFile(path)
	if err != nil {
		err_fatal(err)
	}
	var credentials Credentials
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		err_fatal(err)
	}
	return credentials
}
