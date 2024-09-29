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
}

func create_playlist(db *sql.DB, name string, key string, access_token string) PlaylistResp {
	_, exists := find_row(db, name, "./sql/filter_playlist.sql")
	if exists {
		fmt.Println("playlist already exists")
		os.Exit(0)
	}
	item := create_remote_playlist(name, key, access_token)
	fmt.Println("playlist created!", item.Id, item.Snippet.Title)
	createPlaylistRow(db, item.Id, item.Snippet.Title)
	return item
}

func createPlaylistRow(db *sql.DB, playlist_id string, name string) {
	insertQuery := readSQLFile("./sql/create_playlist.sql")
	_, err := db.Exec(insertQuery, playlist_id, name)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created playlist row")
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
		err = rows.Scan(&playlist.id, &playlist.playlist_id, &playlist.name)
		if err != nil {
			log.Fatal(err)
		}
		playlists = append(playlists, playlist)
	}
	return playlists
}
