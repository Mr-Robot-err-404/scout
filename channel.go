package main

import "database/sql"

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
