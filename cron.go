package main

import (
	"fmt"
	"os"
	"scout/scout_db"
	"strconv"
	"sync"
)

type PlaylistState struct {
	Id    string
	Items int
}

type UpdatedPlaylist struct {
	videos      []Video
	items       int
	playlist_id string
}
type Payload struct {
	state PlaylistState
	err   error
}
type InsertPayload struct {
	updated   UpdatedPlaylist
	err_queue []error
}

func cron_job(api_key string, access_token string, units *int, config Config) ([]UpdatedPlaylist, error) {
	playlists := read_playlists()
	if len(playlists) == 0 {
		info_msg_fatal("no playlists to track")
	}
	updated := []UpdatedPlaylist{}
	route := "https://youtube.googleapis.com/youtube/v3/playlists?part=snippet&maxResults=50&mine=true&key=" + api_key

	log := "fetch remote playlists"
	load(log)
	remote_playlists, err := list_remote_playlist_items(access_token, units, route)
	if err != nil {
		err_msg(log)
		return updated, err
	}
	success_msg(log)
	active, unused := filter_existing_IDs(playlists, remote_playlists)
	err_queue := clear_unused_playlists(unused)
	log_err_queue(err_queue)

	if len(active) == 0 {
		info_msg_fatal("no remote playlists to track")
	}
	var wg sync.WaitGroup
	peek_ch := make(chan Payload)
	done := make(chan bool)

	err_queue = []error{}
	states := []PlaylistState{}
	log = "update playlists"
	load(log)

	for _, playlist_id := range active {
		route := "https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50&playlistId=" + playlist_id + "&key=" + api_key
		wg.Add(1)
		go func() {
			remote_items, err := list_remote_playlist_items(access_token, units, route)
			peek_ch <- Payload{state: PlaylistState{Id: playlist_id, Items: len(remote_items)}, err: err}
			wg.Done()
		}()
	}
	go func() {
		for {
			payload, next := <-peek_ch
			if !next {
				done <- true
				continue
			}
			if payload.err != nil {
				err_queue = append(err_queue, err)
				continue
			}
			states = append(states, payload.state)
		}
	}()
	wg.Wait()
	close(peek_ch)
	<-done

	if len(err_queue) == len(active) || len(states) == 0 {
		err_msg(log)
		log_err_queue(err_queue)
		os.Exit(0)
	}
	log_err_queue(err_queue)
	err_queue = []error{}
	insert_ch := make(chan InsertPayload)

	for _, curr := range states {
		if curr.Items >= config.max_items {
			c := get_item_count(playlists, curr.Id)
			if c > 0 && c != curr.Items {
				updated = append(updated, UpdatedPlaylist{playlist_id: curr.Id, items: curr.Items})
			}
			continue
		}
		wg.Add(1)
		go func() {
			updated_playlist, err_queue := update_playlist(curr.Id, config, units, curr.Items, api_key, access_token)
			insert_ch <- InsertPayload{updated: updated_playlist, err_queue: err_queue}
			wg.Done()
		}()
	}
	go func() {
		for {
			item, next := <-insert_ch
			if !next {
				done <- true
				continue
			}
			if len(item.err_queue) != 0 {
				err_queue = append(err_queue, item.err_queue...)
				continue
			}
			updated = append(updated, item.updated)
		}
	}()
	wg.Wait()
	close(insert_ch)
	<-done

	if len(states) == len(err_queue) {
		err_msg(log)
		log_err_queue(err_queue)
		os.Exit(0)
	}
	success_msg(log)
	log_err_queue(err_queue)
	return updated, nil
}

func get_item_count(playlists []Playlist, playlist_id string) int {
	c := -1
	for _, curr := range playlists {
		if curr.playlist_id != playlist_id {
			continue
		}
		val, err := strconv.Atoi(curr.items)
		if err != nil {
			err_resp(err)
			return c
		}
		c = val
		break
	}
	return c
}

func clear_unused_playlists(IDs []string) []error {
	err_queue := []error{}

	for _, id := range IDs {
		_, err := queries.Delete_playlist(ctx, id)
		if err != nil {
			err_queue = append(err_queue, err)
		}
	}
	return err_queue
}

func update_items(updated []UpdatedPlaylist) {
	for _, curr := range updated {
		params := scout_db.Update_playlist_item_count_params{Items: int64(curr.items), PlaylistID: curr.playlist_id}
		err := queries.Update_playlist_item_count(ctx, params)
		if err != nil {
			err_resp(err)
			continue
		}
	}

}

func update_vids(updated []UpdatedPlaylist) {
	for _, curr := range updated {
		add_vid_rows(curr.videos)
	}
}

func update_playlist(playlist_id string, config Config, units *int, items int, api_key string, access_token string) (UpdatedPlaylist, []error) {
	updated := UpdatedPlaylist{}
	playlist, err := queries.Find_playlist(ctx, playlist_id)
	if err != nil {
		return updated, []error{err}
	}
	q := parse_input(playlist.Q)
	filter := parse_input(playlist.Filter)

	limit := config.max_items - items
	playlist_items, err_queue := select_playlist_items(q, filter, units, limit, config.format, config.category)
	if len(err_queue) != 0 {
		return updated, err_queue
	}
	if len(playlist_items) == 0 {
		err = fmt.Errorf("no matching search results for ID: %v", playlist_id)
		return updated, []error{err}
	}
	videos, c, err_queue := populate_playlist(playlist_id, units, playlist_items, api_key, access_token)
	updated.items = c + items
	updated.videos = videos
	updated.playlist_id = playlist_id

	return updated, err_queue
}

func filter_existing_IDs(playlists []Playlist, existing []string) ([]string, []string) {
	valid := []string{}
	remove := []string{}

	for i := range playlists {
		curr := playlists[i]
		exists := false
		for _, id := range existing {
			if curr.playlist_id == id {
				exists = true
				break
			}
		}
		if exists {
			valid = append(valid, curr.playlist_id)
			continue
		}
		remove = append(remove, curr.playlist_id)
	}
	return valid, remove
}
