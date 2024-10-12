package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Snippet struct {
	Title     string
	CustomUrl string
}

type Item struct {
	Id      string
	Snippet Snippet
}
type ChannelResp struct {
	Items []Item
}
type PlaylistResp struct {
	Id string
	PlaylistSnippet
}
type PlaylistSnippet struct {
	Snippet struct {
		Title string `json:"title"`
	} `json:"snippet"`
}

type PlaylistInfo struct {
	Snippet PlaylistInfoSnippet `json:"snippet"`
}

type PlaylistInfoSnippet struct {
	PlaylistID string     `json:"playlistId"`
	ResourceID ResourceID `json:"resourceId"`
}

type PlaylistInsertResp struct {
	Snippet SearchSnippet
}

type ResourceID struct {
	VideoID string `json:"videoId"`
	Kind    string `json:"kind"`
}

type SearchSnippet struct {
	ChannelId    string
	Title        string
	Description  string
	ChannelTitle string
}

type SearchItem struct {
	Id      SearchID
	Snippet SearchSnippet
}

type SearchID struct {
	Kind       string
	VideoId    string
	PlaylistId string
}

type SearchResp struct {
	Items    []SearchItem
	PageInfo PageInfo
}

type PageInfo struct {
	TotalResults   int
	ResultsPerPage int
}

func create_remote_playlist(playlist_name string, key string, access_token string) (PlaylistResp, error) {
	var item PlaylistResp
	var snippet PlaylistSnippet
	snippet.Snippet.Title = playlist_name

	route := "https://youtube.googleapis.com/youtube/v3/playlists?part=snippet&videoDuration=long&key=" + key
	json_body, err := json.Marshal(&snippet)
	if err != nil {
		return item, err
	}
	body_reader := bytes.NewReader(json_body)

	req, err := http.NewRequest(http.MethodPost, route, body_reader)
	if err != nil {
		return item, err
	}
	req.Header.Set("Authorization", "Bearer "+access_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return item, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("failed request with status: %v", resp.Status)
		return item, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return item, err
	}
	err = json.Unmarshal(body, &item)

	if err != nil {
		return item, err
	}
	return item, nil
}

func search_remote_channel(q string, channel_id string) (SearchResp, error) {
	var res SearchResp
	api_key := os.Getenv("API_KEY")

	route := "https://youtube.googleapis.com/youtube/v3/search?part=snippet&channelId=" + channel_id + "&type=video&q=" + q + "&key=" + api_key
	resp, err := http.Get(route)

	if err != nil {
		return res, err
	}
	if resp.StatusCode != 200 {
		return res, fmt.Errorf("request failed with status code %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return res, err
	}
	err = json.Unmarshal(body, &res)

	if err != nil {
		return res, err
	}
	return res, nil
}

func insert_playlist_item(playlist_id string, video_id string, key string, access_token string) (PlaylistInsertResp, error) {
	var item PlaylistInfo
	item.Snippet.PlaylistID = playlist_id
	item.Snippet.ResourceID.VideoID = video_id
	item.Snippet.ResourceID.Kind = "youtube#video"

	route := "https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&key=" + key

	json_body, err := json.Marshal(&item)
	if err != nil {
		return PlaylistInsertResp{}, err
	}
	body_reader := bytes.NewReader(json_body)

	req, err := http.NewRequest(http.MethodPost, route, body_reader)
	if err != nil {
		return PlaylistInsertResp{}, err
	}
	req.Header.Set("Authorization", "Bearer "+access_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PlaylistInsertResp{}, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("failed request with status: %v", resp.StatusCode)
		return PlaylistInsertResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PlaylistInsertResp{}, err
	}
	var res PlaylistInsertResp
	err = json.Unmarshal(body, &res)

	if err != nil {
		return PlaylistInsertResp{}, err
	}
	return res, nil
}

func get_channel_ID(tag string, key string) ([]string, error) {
	term := tag
	if string(term[0]) == "@" {
		term = term[1:]
	}
	q := parse_query(term)
	route := "https://youtube.googleapis.com/youtube/v3/channels?part=snippet&forHandle=%40" + q + "&key=" + key
	resp, err := http.Get(route)

	if err != nil {
		return []string{}, err
	}
	if resp.StatusCode != 200 {
		return []string{}, fmt.Errorf("request failed with status %v", resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []string{}, err
	}
	var s ChannelResp
	err = json.Unmarshal(body, &s)

	if err != nil {
		return []string{}, err
	}
	if len(s.Items) == 0 {
		return []string{}, fmt.Errorf("no channels were found for %s", tag)
	}
	return []string{s.Items[0].Id, s.Items[0].Snippet.Title, s.Items[0].Snippet.CustomUrl}, nil
}
