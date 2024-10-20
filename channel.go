package main

import "database/sql"

type Channel struct {
	channel_id string
	tag        string
	name       string
	category   string
}

func read_channels(db *sql.DB) []Channel {
	channels := []Channel{}
	query := readSQLFile("./sql/read_all_channels.sql")
	rows, err := db.Query(query)
	if err != nil {
		err_fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var channel Channel
		err = rows.Scan(&channel.channel_id, &channel.tag, &channel.name, &channel.category)
		if err != nil {
			err_fatal(err)
		}
		channels = append(channels, channel)
	}
	return channels
}

func drop_channel_table(db *sql.DB) {
	query := "DROP TABLE channel"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table")
}
