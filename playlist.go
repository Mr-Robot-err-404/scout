package main

import (
	"database/sql"
	"fmt"
	"os"
)

type Playlist struct {
	playlist_id string
	name        string
	q           string
	filter      string
	long_format string
}

func create_playlist(db *sql.DB, name string, q string, filter string, key string, access_token string) PlaylistResp {
	_, exists := find_row(db, name, "./sql/filter_playlist.sql")
	if exists {
		info_msg_fatal("playlist already exists")
	}
	log := "create remote playlist"
	load(log)
	item, err := create_remote_playlist(name, key, access_token)
	if err != nil {
		err_msg(log)
		err_fatal(err)
	}
	success_msg(log)

	create_playlist_row(db, item.Id, item.Snippet.Title, q, filter, true)
	return item
}

func populate_playlist(db *sql.DB, q []string, filter []string, playlist_id string) []Video {
	search_items := [][]SearchItem{}
	channels := readChannels(db)
	query := convert_and_parse(q)

	for i := range channels {
		curr := channels[i]
		resp, err := search_remote_channel(query, curr.channel_id)
		if err != nil {
			err_msg(curr.name)
			continue
		}
		success_msg(curr.name)
		search_items = append(search_items, resp.Items)
	}
	playlist_items := rank_channels(search_items, filter)

	if len(playlist_items) == 0 {
		info_msg_fatal("no matching search results")
	}
	api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
	videos := []Video{}

	c := 0
	for i := range playlist_items {
		video_id := playlist_items[i].Id.VideoId
		if len(video_id) == 0 {
			continue
		}
		_, err := insert_playlist_item(playlist_id, video_id, api_key, access_token)
		if err != nil {
			err_resp(err)
			continue
		}
		vid := Video{title: playlist_items[i].Snippet.Title, video_id: video_id}
		videos = append(videos, vid)
		c++
	}
	msg := fmt.Sprintf("added %v items to playlist\n", c)
	info_msg(msg)

	return videos
}

func create_playlist_row(db *sql.DB, playlist_id string, name string, q string, filter string, long bool) {
	insertQuery := readSQLFile("./sql/create_playlist.sql")
	_, err := db.Exec(insertQuery, playlist_id, name, q, filter, long)
	if err != nil {
		err_fatal(err)
	}
}

func delete_playlist(db *sql.DB, name string) error {
	query := readSQLFile("./sql/delete_playlist_row.sql")
	_, err := db.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func read_playlists(db *sql.DB) []Playlist {
	playlists := []Playlist{}

	query := readSQLFile("./sql/read_all_playlists.sql")
	rows, err := db.Query(query)
	if err != nil {
		err_fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var playlist Playlist
		err = rows.Scan(&playlist.playlist_id, &playlist.name, &playlist.q, &playlist.filter, &playlist.long_format)
		if err != nil {
			err_fatal(err)
		}
		playlists = append(playlists, playlist)
	}
	return playlists
}

func drop_playlist_table(db *sql.DB) {
	query := "DROP TABLE playlist"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table")
}
