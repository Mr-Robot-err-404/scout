package main

import (
	"fmt"
	"scout/scout_db"
	"strings"
)

type PlaylistFields struct {
	q        []string
	filter   []string
	category string
	format   string
}

func edit_playlist(playlist_id string) error {
	playlist, err := queries.Find_playlist_id(ctx, playlist_id)
	if err != nil {
		return err
	}
	print_title_with_bg("edit " + playlist.Name)
	next_q, err := edit_user_input("Search query: ", playlist.Q)

	if err != nil {
		err_fatal(err)
	}
	next_filter, err := edit_user_input("Filter: ", playlist.Filter)

	if err != nil {
		err_fatal(err)
	}
	next_category, err := edit_user_input("Category: ", playlist.Category)

	if err != nil {
		err_fatal(err)
	}
	next_format, err := edit_user_input("Format: ", playlist.Format)

	if err != nil {
		err_fatal(err)
	}
	original := get_playlist_fields(playlist.Q, playlist.Filter, playlist.Category, playlist.Format)
	next := get_playlist_fields(next_q, next_filter, next_category, next_format)

	if !is_playlist_edited(original, next) {
		return nil
	}
	params, err := get_edited_fields(original, next, playlist_id)
	if err != nil {
		return err
	}
	err = queries.Update_playlist(ctx, params)
	if err != nil {
		return err
	}
	success_msg("edited playlist")
	return nil
}

func get_edited_fields(original PlaylistFields, next PlaylistFields, playlist_id string) (scout_db.Update_playlist_params, error) {
	csv_q := csv_string(original.q)
	csv_filter := csv_string(original.filter)
	params := scout_db.Update_playlist_params{Q: csv_q, Filter: csv_filter, Category: original.category, PlaylistID: playlist_id}

	if is_edited(original.q, next.q) {
		csv := csv_string(next.q)
		params.Q = csv
	}
	if is_edited(original.filter, next.filter) {
		csv := csv_string(next.filter)
		params.Filter = csv
	}
	if original.category != next.category {
		params.Category = next.category
	}
	if original.format != next.format {
		format := strings.TrimSpace(next.format)
		if !is_format_valid(format) {
			err := fmt.Errorf("invalid video format. Accepted values: short || medium || long")
			return params, err
		}
		params.Format = format
	}
	return params, nil
}

func is_playlist_edited(original PlaylistFields, next PlaylistFields) bool {
	if is_edited(original.q, next.q) {
		return true
	}
	if is_edited(original.filter, next.filter) {
		return true
	}
	if original.category != next.category {
		return true
	}
	if original.format != next.format {
		return true
	}
	return false
}

func get_playlist_fields(q string, filter string, category string, format string) PlaylistFields {
	var fields PlaylistFields
	fields.q = parse_input(q)
	fields.filter = parse_input(filter)
	fields.category = category
	fields.format = format
	return fields
}

func is_format_valid(format string) bool {
	if format == "long" || format == "medium" || format == "short" {
		return true
	}
	return false
}

func is_edited(original []string, input []string) bool {
	if len(original) != len(input) {
		return true
	}
	for i := range original {
		if original[i] != input[i] {
			return true
		}
	}
	return false
}
