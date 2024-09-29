package main

import (
	"database/sql"
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
