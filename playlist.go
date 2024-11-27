package main

import (
	"fmt"
	"sync"
)

type SearchPayload struct {
	resp SearchResp
	err  error
}

type Playlist struct {
	playlist_id string
	name        string
	q           string
	filter      string
	format      string
	items       string
	category    string
}

func create_playlist(name string, key string, access_token string, units *int) PlaylistResp {
	log := "create remote playlist"
	load(log)
	item, err := create_remote_playlist(name, key, access_token)
	if err != nil {
		err_msg(log)
		err_fatal(err)
	}
	success_msg(log)
	*units -= quota_map["insert"]

	return item
}

func populate_playlist(playlist_id string, units *int, playlist_items []SearchItem, api_key string, access_token string) ([]Video, int, []error) {
	videos := []Video{}
	err_queue := []error{}
	c := 0

	for i := range playlist_items {
		video_id := playlist_items[i].Id.VideoId
		if len(video_id) == 0 {
			continue
		}
		_, err := insert_playlist_item(playlist_id, video_id, api_key, access_token)
		if err != nil {
			err_queue = append(err_queue, err)
			continue
		}
		vid := Video{title: playlist_items[i].Snippet.Title, video_id: video_id}
		videos = append(videos, vid)
		c++
		*units -= quota_map["insert"]
	}
	return videos, c, err_queue
}

func select_playlist_items(q []string, filter []string, units *int, max_items int, format string, category string) ([]SearchItem, []error) {
	search_items := [][]SearchItem{}
	query := convert_and_parse(q)
	yt_channels, err := queries.Channels_by_category(ctx, category)
	if err != nil {
		return []SearchItem{}, []error{err}
	}
	vids, err := read_videos()
	if err != nil {
		return []SearchItem{}, []error{err}
	}
	if len(yt_channels) == 0 {
		msg := fmt.Sprintf("no channels tracked for category: %v", category)
		info_msg_fatal(msg)
	}
	var wg sync.WaitGroup
	ch := make(chan SearchPayload)
	done := make(chan bool)

	err_queue := []error{}

	for _, curr := range yt_channels {
		wg.Add(1)
		go search_remote_channel(query, curr.ChannelID, format, &wg, ch)
	}
	go func() {
		for {
			item, next := <-ch
			if !next {
				done <- true
				continue
			}
			if item.err != nil {
				err_queue = append(err_queue, err)
				continue
			}
			search_items = append(search_items, item.resp.Items)
			*units -= quota_map["search"]
		}
	}()
	wg.Wait()
	close(ch)
	<-done

	if len(err_queue) == len(yt_channels) {
		return []SearchItem{}, err_queue
	}
	playlist_items := rank_channels(search_items, filter, vids, max_items)
	return playlist_items, nil
}

func read_playlists() []Playlist {
	playlists := []Playlist{}

	query := "SELECT * FROM playlist;"
	rows, err := db.Query(query)
	if err != nil {
		err_fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var playlist Playlist
		err = rows.Scan(&playlist.playlist_id, &playlist.name, &playlist.q, &playlist.filter, &playlist.format, &playlist.items, &playlist.category)
		if err != nil {
			err_fatal(err)
		}
		playlists = append(playlists, playlist)
	}
	return playlists
}

func clear_playlist_table() {
	query := "DELETE FROM playlist"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("cleared table")
}

func drop_playlist_table() {
	query := "DROP TABLE playlist"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table")
}
