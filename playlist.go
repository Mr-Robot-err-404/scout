package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Playlist struct {
	id               string
	playlist_id      string
	name             string
	q                string
	inclusive_search bool
	long_format      string
}

func create_playlist(db *sql.DB, name string, q string, key string, access_token string) PlaylistResp {
	_, exists := find_row(db, name, "./sql/filter_playlist.sql")
	if exists {
		fmt.Println("playlist already exists")
		os.Exit(0)
	}
	item := create_remote_playlist(name, key, access_token)
	fmt.Println("playlist created!", item.Id, item.Snippet.Title)

	createPlaylistRow(db, item.Id, item.Snippet.Title, q, true, true)
	return item
}

func createPlaylistRow(db *sql.DB, playlist_id string, name string, q string, inclusive bool, long bool) {
	insertQuery := readSQLFile("./sql/create_playlist.sql")
	_, err := db.Exec(insertQuery, playlist_id, name, q, inclusive, long)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created playlist row")
}

func deletePlaylist(db *sql.DB, name string) error {
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
		err = rows.Scan(&playlist.id, &playlist.playlist_id, &playlist.name, &playlist.q, &playlist.inclusive_search, &playlist.long_format)
		if err != nil {
			log.Fatal(err)
		}
		playlists = append(playlists, playlist)
	}
	return playlists
}
