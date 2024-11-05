package main

import (
	"fmt"
	"scout/scout_db"
	"strconv"
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

func cron_job(api_key string, access_token string, units *int) ([]UpdatedPlaylist, error) {
	updated := []UpdatedPlaylist{}
	playlists := read_playlists()
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
	clear_unused_playlists(unused)

	config, err := read_config_file()
	if err != nil {
		return updated, err
	}

	states := []PlaylistState{}
	for _, playlist_id := range active {
		route := "https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50&playlistId=" + playlist_id + "&key=" + api_key
		remote_items, err := list_remote_playlist_items(access_token, units, route)
		if err != nil {
			err_resp(err)
			continue
		}
		state := PlaylistState{Id: playlist_id, Items: len(remote_items)}
		states = append(states, state)
	}
	for _, curr := range states {
		if curr.Items >= config.max_items {
			c := get_item_count(playlists, curr.Id)
			if c > 0 && c != curr.Items {
				updated = append(updated, UpdatedPlaylist{playlist_id: curr.Id, items: curr.Items})
			}
			continue
		}
		updated_playlist, err := update_playlist(curr.Id, config, units, curr.Items, api_key, access_token)
		if err != nil {
			err_msg(err.Error())
			continue
		}
		updated = append(updated, updated_playlist)
	}
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

func clear_unused_playlists(IDs []string) {
	// TODO: delete batch

	for _, id := range IDs {
		_, err := queries.Delete_playlist(ctx, id)
		if err != nil {
			msg := fmt.Sprintf("failed to delete record with playlist_id: %v", id)
			err_msg(msg)
		}
	}
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

func update_playlist(playlist_id string, config Config, units *int, items int, api_key string, access_token string) (UpdatedPlaylist, error) {
	updated := UpdatedPlaylist{}
	playlist, err := queries.Find_playlist(ctx, playlist_id)
	if err != nil {
		return updated, err
	}
	q := parse_input(playlist.Q)
	filter := parse_input(playlist.Filter)

	limit := config.max_items - items
	playlist_items, err := select_playlist_items(q, filter, units, limit, config.format, config.category)
	if err != nil {
		return updated, err
	}
	if len(playlist_items) == 0 {
		msg := fmt.Sprintf("no matching search results for ID: %v", playlist_id)
		info_msg_fatal(msg)
	}
	videos, c := populate_playlist(playlist_id, units, playlist_items, api_key, access_token)
	updated.items = c + items
	updated.videos = videos
	updated.playlist_id = playlist_id

	return updated, nil
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
