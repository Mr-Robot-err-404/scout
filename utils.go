package main

import (
	"database/sql"
	"net/url"
	"strings"
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
func get_playlist_display(playlists []Playlist) ([]string, [][]string) {
	headers := []string{"NAME", "QUERY", "FILTER", "ITEMS", "VID_FORMAT", "PLAYLIST_ID"}
	display_rows := [][]string{}

	for i := range playlists {
		curr := playlists[i]
		q := get_display_query(strings.Split(curr.q, ","))
		filter := get_display_query(strings.Split(curr.filter, ","))
		row := []string{curr.name, q, filter, "0", curr.long_format, curr.playlist_id}
		display_rows = append(display_rows, row)
	}
	return headers, display_rows
}

func get_display_query(query []string) string {
	if len(query) == 0 {
		return ""
	}
	q := query[0]
	for i := 1; i < len(query); i++ {
		q += " || " + query[i]
	}
	return q
}

func get_channel_display(channels []Channel) ([]string, [][]string) {
	headers := []string{"NAME", "TAG", "CATEGORY", "CHANNEL_ID"}
	display_rows := [][]string{}

	for i := range len(channels) {
		curr := channels[i]
		row := []string{curr.name, curr.tag, curr.category, curr.channel_id}
		display_rows = append(display_rows, row)
	}
	return headers, display_rows
}

func get_channel_IDs(channels []Channel) []string {
	IDs := []string{}
	for i := range channels {
		curr := channels[i]
		IDs = append(IDs, curr.channel_id)
	}
	return IDs
}
