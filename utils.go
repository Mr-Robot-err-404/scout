package main

import (
	"database/sql"
	"fmt"
	"net/url"
)

func find_row(db *sql.DB, search_term string, path string) (string, bool) {
	var tag string
	s := search_term[:]

	if string(s[0]) != "@" {
		s = "@" + s
	}
	query := readSQLFile(path)
	row := db.QueryRow(query, s)
	err := row.Scan(&tag)
	if err != nil {
		return "", false
	}
	return tag, true
}

func csv_string(q []string) string {
	csv_line := ""
	for i := range q {
		str := q[i]
		if i == 0 {
			csv_line += str
			continue
		}
		csv_line += "," + str
	}
	return csv_line
}

func parse_query(q string) string {
	return url.PathEscape(q)
}

func convert_and_parse(q []string) string {
	if len(q) == 0 {
		return ""
	}
	query := q[0]
	if len(q) == 1 {
		return parse_query(query)
	}
	for i := 1; i < len(q); i++ {
		s := q[i]
		query += "|"
		query += s
	}
	return parse_query(query)
}

func show_gcloud_tokens(access_token string) {
	credentials := readCredentialsFile("../.config/gcloud/application_default_credentials.json")
	fmt.Println("----------------------------------------------")
	fmt.Printf("REFRESH_TOKEN %v\n", credentials.Refresh_token)
	fmt.Println("----------------------------------------------")
	fmt.Printf("CLIENT_ID     %v\n", credentials.Client_id)
	fmt.Println("----------------------------------------------")
	fmt.Printf("CLIENT_SECRET %v\n", credentials.Client_secret)
	fmt.Println("----------------------------------------------")
	fmt.Printf("ACCESS_TOKEN %v\n", access_token)
	fmt.Println("----------------------------------------------")

}

func get_channel_IDs(channels []Channel) []string {
	IDs := []string{}
	for i := range channels {
		curr := channels[i]
		IDs = append(IDs, curr.channel_id)
	}
	return IDs
}
