package main

import "strings"

func get_playlist_display(playlists []Playlist) ([]string, [][]string) {
	headers := []string{"NAME", "QUERY", "FILTER", "ITEMS", "VID_FORMAT", "PLAYLIST_ID"}
	display_rows := [][]string{}

	for i := range playlists {
		curr := playlists[i]
		q := get_display_query(strings.Split(curr.q, ","))
		filter := get_display_query(strings.Split(curr.filter, ","))
		row := []string{curr.name, q, filter, curr.items, curr.format, curr.playlist_id}
		display_rows = append(display_rows, row)
	}
	return headers, display_rows
}

func get_display_query(query []string) string {
	if len(query) == 0 {
		return ""
	}
	q := query[0]
	for i := 1; i < len(query); i++ {
		q += " || " + query[i]
	}
	return q
}

func get_channel_display(channels []Channel) ([]string, [][]string) {
	headers := []string{"NAME", "TAG", "CATEGORY", "CHANNEL_ID"}
	display_rows := [][]string{}

	for i := range len(channels) {
		curr := channels[i]
		row := []string{curr.name, curr.tag, curr.category, curr.channel_id}
		display_rows = append(display_rows, row)
	}
	return headers, display_rows
}

func get_video_display(videos []Video) ([]string, [][]string) {
	headers := []string{"VIDEO_ID", "TITLE"}
	display_rows := [][]string{}

	for i := range len(videos) {
		curr := videos[i]
		row := []string{curr.video_id, curr.title}
		display_rows = append(display_rows, row)
	}
	return headers, display_rows
}
