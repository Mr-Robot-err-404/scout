package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Quota struct {
	id           int
	quota        int
	quota_reset  time.Time
	last_refresh time.Time
}
type OAuth struct {
	Access_token string
}

func refresh_token() (string, error) {
	refresh_token := os.Getenv("REFRESH_TOKEN")
	client_id, client_secret := os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")

	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("refresh_token", refresh_token)
	data.Set("grant_type", "refresh_token")

	route := "https://www.googleapis.com/oauth2/v4/token"
	resp, err := http.PostForm(route, data)

	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf(resp.Status)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	var s OAuth
	err = json.Unmarshal(body, &s)

	if err != nil {
		return "", err
	}
	err = queries.Update_last_refresh(ctx)
	if err != nil {
		return "", err
	}
	renew_access_token(s.Access_token)
	return s.Access_token, nil
}

func renew_access_token(access_token string) error {
	env, err := get_env_map()
	if err != nil {
		return err
	}
	env["ACCESS_TOKEN"] = access_token
	err = godotenv.Write(env, "./.env")
	if err != nil {
		return err
	}
	return nil
}

func get_env_map() (map[string]string, error) {
	var env_map map[string]string
	env_map, err := godotenv.Read()

	if err != nil {
		return env_map, err
	}
	return env_map, nil
}

func check_token() string {
	quota, err := read_quota()
	if err != nil {
		err_fatal(err)
	}
	ts := quota.last_refresh
	diff := time.Now().Unix() - ts.Unix()

	if diff < 3000 {
		return ""
	}
	log := "refresh access token"
	load(log)

	token, err := refresh_token()
	if err != nil {
		err_msg(log)
		err_fatal(err)
	}
	success_msg(log)
	return token
}
