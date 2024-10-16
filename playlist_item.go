package main

import (
	"database/sql"
	"strconv"
)

type PlaylistItem struct {
	id          int
	video_id    string
	playlist_id string
	channel_id  string
}

func insert_item_row(db *sql.DB, video_id string, playlist_id string, channel_id string) error {
	query := readSQLFile("./sql/create_playlist_item.sql")
	_, err := db.Exec(query, video_id, playlist_id, channel_id)
	if err != nil {
		return err
	}
	return nil
}

func read_playlist_items(db *sql.DB) ([]PlaylistItem, error) {
	items := []PlaylistItem{}
	query := "SELECT * FROM playlist_item"
	rows, err := db.Query(query)
	if err != nil {
		return items, err
	}
	defer rows.Close()
	for rows.Next() {
		var item PlaylistItem
		err = rows.Scan(&item.id, &item.video_id, &item.playlist_id, &item.channel_id)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}
	return items, nil
}

func get_items_display(items []PlaylistItem) ([]string, [][]string) {
	headers := []string{"ID", "VIDEO_ID", "PLAYLIST_ID", "CHANNEL_ID"}
	display_rows := [][]string{}

	for i := range items {
		curr := items[i]
		id := strconv.Itoa(curr.id)
		row := []string{id, curr.video_id, curr.playlist_id, curr.channel_id}
		display_rows = append(display_rows, row)
	}
	return headers, display_rows
}

func drop_item_table(db *sql.DB) {
	query := "DROP TABLE playlist_item"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table")
}

func clear_item_table(db *sql.DB) {
	query := "DELETE FROM playlist_item"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("cleared table playlist_item")

}