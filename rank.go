package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func rank_channels(lists [][]SearchItem, q []string) []SearchItem {
	filtered_lists := [][]SearchItem{}

	for i := range lists {
		items := lists[i]
		filter_list(&items, q)
		filtered_lists = append(filtered_lists, items)
	}
	max := max_len(filtered_lists)
	items := fill_playlist(filtered_lists, max, 5)
	return items
}

func fill_playlist(lists [][]SearchItem, max int, capacity int) []SearchItem {
	playlist_items := []SearchItem{}

	for i := 0; i < max; i++ {
		for j := 0; j < len(lists); j++ {
			if len(playlist_items) == capacity {
				return playlist_items
			}
			item := lists[j]
			if i >= len(item) {
				continue
			}
			playlist_items = append(playlist_items, item[i])
		}
	}
	return playlist_items
}

func filter_list(videos *[]SearchItem, q []string) {
	items := *videos
	video_items := []SearchItem{}

	for i := range items {
		curr := items[i]
		title, desc := curr.Snippet.Title, curr.Snippet.Description
		valid := 0

		if len(title) == 0 || len(desc) == 0 {
			continue
		}
		for _, s := range q {
			valid += matching_terms(s, strings.ToLower(title), strings.ToLower(desc))
		}
		if valid == 0 {
			continue
		}
		video_items = append(video_items, curr)
	}
	*videos = video_items
}

func matching_terms(q string, title string, desc string) int {
	slice := strings.Split(q, ",")
	is_title_match := false
	is_desc_match := false

	for i := range slice {
		str := slice[i]
		if is_substring(title, strings.Split(str, " ")) {
			is_title_match = true
		}
		if is_substring(desc, strings.Split(str, " ")) {
			is_desc_match = true
		}
	}
	if is_title_match && is_desc_match {
		return 2
	}
	if is_title_match {
		return 1
	}
	if is_desc_match {
		return 1
	}
	return 0
}

func is_substring(str string, sub []string) bool {
	for _, s := range sub {
		is_substring := strings.Contains(str, s)
		if !is_substring {
			return false
		}
	}
	return true
}

func max_len(lists [][]SearchItem) int {
	max := -1
	for _, curr := range lists {
		if len(curr) > max {
			max = len(curr)
		}
	}
	return max
}

func read_sample(path string) SearchResp {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var res SearchResp
	err = json.Unmarshal(data, &res)

	if err != nil {
		log.Fatal(err)
	}
	return res
}
