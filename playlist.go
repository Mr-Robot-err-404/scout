package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Playlist struct {
	id          string
	playlist_id string
	name        string
	q           string
	filter      string
	long_format string
}

func create_playlist(db *sql.DB, name string, q string, filter string, key string, access_token string) PlaylistResp {
	_, exists := find_row(db, name, "./sql/filter_playlist.sql")
	if exists {
		fmt.Println("playlist already exists")
		os.Exit(0)
	}
	item := create_remote_playlist(name, key, access_token)
	fmt.Println("created playlist")

	create_playlist_row(db, item.Id, item.Snippet.Title, q, filter, true)
	return item
}

func populate_playlist(db *sql.DB, q []string, filter []string, playlist_id string) {
	search_items := [][]SearchItem{}
	channels := readChannels(db)
	query := convert_and_parse(q)

	fmt.Println("searching channels...")

	for i := range channels {
		curr := channels[i]
		resp, err := search_remote_channel(query, curr.channel_id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		search_items = append(search_items, resp.Items)
	}
	playlist_items := rank_channels(search_items, filter)

	if len(playlist_items) == 0 {
		fmt.Println("no items found")
		os.Exit(0)
	}
	api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
	fmt.Println("inserting items into playlist...")

	c := 0
	for i := range playlist_items {
		video_id := playlist_items[i].Id.VideoId
		if len(video_id) == 0 {
			continue
		}
		_, err := insert_playlist_item(playlist_id, video_id, api_key, access_token)
		if err != nil {
			fmt.Println(err)
			continue
		}
		c++
	}
	fmt.Printf("added %v items to playlist\n", c)
}

func create_playlist_row(db *sql.DB, playlist_id string, name string, q string, filter string, long bool) {
	insertQuery := readSQLFile("./sql/create_playlist.sql")
	_, err := db.Exec(insertQuery, playlist_id, name, q, filter, long)
	if err != nil {
		log.Fatal(err)
	}
}

func delete_playlist(db *sql.DB, name string) error {
	query := readSQLFile("./sql/delete_playlist_row.sql")
	_, err := db.Exec(query, name)
	if err != nil {
		log.Printf("error deleting table: %v", err)
		return err
	}
	fmt.Println("playlist removed")
	return nil
}

func read_playlists(db *sql.DB) []Playlist {
	playlists := []Playlist{}

	query := readSQLFile("./sql/read_all_playlists.sql")
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var playlist Playlist
		err = rows.Scan(&playlist.id, &playlist.playlist_id, &playlist.name, &playlist.q, &playlist.filter, &playlist.long_format)
		if err != nil {
			log.Fatal(err)
		}
		playlists = append(playlists, playlist)
	}
	return playlists
}
