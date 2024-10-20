package main

import (
	"database/sql"
)

func init_quota_map() map[string]int {
	quota_map := map[string]int{
		"get":    1,
		"search": 50,
		"insert": 100,
	}
	return quota_map
}

func read_quota(db *sql.DB) (Quota, error) {
	var quota Quota
	query := "SELECT * FROM quota"
	rows, err := db.Query(query)
	if err != nil {
		return quota, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&quota.id, &quota.quota, &quota.timestamp, &quota.last_refresh)
		if err != nil {
			return quota, err
		}
	}
	return quota, nil
}

func update_quota(db *sql.DB, units *int) {
	query := readSQLFile("./sql/update_quota.sql")
	_, err := db.Exec(query, *units)
	if err != nil {
		err_fatal(err)
	}
}

func init_quota_row(db *sql.DB) {
	query := readSQLFile("./sql/init_quota.sql")
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("quota table initialized")
}

func drop_quota_table(db *sql.DB) {
	query := "DROP TABLE quota"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table quota")
}
