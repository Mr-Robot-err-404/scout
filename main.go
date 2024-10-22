package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var sp = create_spinner()
var quota_map = init_quota_map()

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		help_txt, err := os.ReadFile("./help.txt")
		if err != nil {
			err_fatal(err)
		}
		fmt.Print(string(help_txt))
		return
	}
	err := godotenv.Load(".env")
	if err != nil {
		err_fatal(err)
	}
	db := setup_db()
	defer db.Close()

	cmd := flag.NewFlagSet("create_cmd", flag.ExitOnError)
	create_flag := cmd.String("add", "", "add")
	delete_flag := cmd.String("delete", "", "delete")
	edit_flag := cmd.String("edit", "", "edit")

	config_cmd := flag.NewFlagSet("config_cmd", flag.ExitOnError)
	format_flag := config_cmd.String("format", "", "format")
	category_flag := config_cmd.String("category", "", "category")
	max_flag := config_cmd.String("max", "", "max")

	// TODO: edit cmd for playlists & channel

	switch os.Args[1] {
	case "cli":
		err := insert_item_row(db, "If_5JtqQ4y0", "PL-vGMW-bu9eVAQ-J8GTtLh4DLO4BtvVkg", "UCXy10-NEFGxQ3b4NVrzHw1Q")
		if err != nil {
			err_fatal(err)
		}
		success_msg("inserted row")
	case "vid":
		videos, err := read_videos(db)
		if err != nil {
			err_fatal(err)
		}
		headers, display_rows := get_video_display(videos)
		print_table(headers, display_rows)

	case "channel":
		if len(os.Args) == 2 {
			channels := read_channels(db)
			headers, display_rows := get_channel_display(channels)
			print_table(headers, display_rows)
			return
		}
		cmd.Parse(os.Args[2:])
		if len(*delete_flag) != 0 {
			tag, exists := find_row(db, *delete_flag, "./sql/read_row.sql")
			if !exists {
				err_msg("no channel found with that tag")
				return
			}
			log := "delete channel"
			err = deleteRow(db, tag)
			if err != nil {
				err_msg(log)
				err_fatal(err)
			}
			success_msg(log)
			return
		}
		if len(*create_flag) == 0 {
			err_msg("no channel tag provided")
			return
		}
		_, exists := find_row(db, *create_flag, "./sql/read_row.sql")
		if exists {
			info_msg_fatal("channel is already tracked")
		}
		log := "add channel"
		load(log)
		quota, err := read_quota(db)

		if err != nil {
			err_fatal(err)
		}
		units := quota.quota
		defer update_quota(db, &units)

		key := os.Getenv("API_KEY")
		item, err := get_channel_ID(*create_flag, key, &units)
		if err != nil {
			err_msg(log)
			err_fatal(err)
		}
		id, title, real_tag := item[0], item[1], item[2]
		err = create_channel_row(db, id, real_tag, title)

		if err != nil {
			err_msg(log)
			err_fatal(err)
		}
		success_msg(log)

	case "play":
		if len(os.Args) == 2 {
			playlists := read_playlists(db)
			headers, display_rows := get_playlist_display(playlists)
			print_table(headers, display_rows)
			return
		}
		cmd.Parse(os.Args[2:])

		if len(*delete_flag) != 0 {
			err := delete_playlist(db, *delete_flag)
			if err != nil {
				err_fatal(err)
			}
			success_msg("playlist deleted")
			return
		}
		if len(*edit_flag) != 0 {
			// TODO: edit looks for playlist_id
		}
		query := get_user_input("Enter search terms: ", true)
		filter := get_user_input("Filter: ", false)

		quota, err := read_quota(db)
		if err != nil {
			err_fatal(err)
		}
		units := quota.quota
		defer update_quota(db, &units)

		check_token(db)
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		config, err := read_config_file()
		if err != nil {
			err_fatal(err)
		}
		q := csv_string(query)
		f := csv_string(filter)

		playlist_resp := create_playlist(db, *create_flag, api_key, access_token, &units)
		videos, c := populate_playlist(db, query, filter, playlist_resp.Id, &units, config.max_items)
		add_playlist_row(db, playlist_resp.Id, playlist_resp.Snippet.Title, q, f, c, config.format, config.category)
		add_vid_rows(db, videos)

	case "config":
		config, err := read_config_file()
		if err != nil {
			err_fatal(err)
		}
		if len(os.Args) == 2 {
			print_config(config)
			return
		}
		config_cmd.Parse(os.Args[2:])
		config_options, err := validate_config_flags(*format_flag, *category_flag, *max_flag)
		if err != nil {
			err_fatal(err)
		}
		err = update_config_file(config_options)
		if err != nil {
			err_fatal(err)
		}
		success_msg("updated config file")

	case "table":
		err := createTable(db, "./sql/create_playlist_table.sql")
		if err != nil {
			err_fatal(err)
		}
		success_msg("created table")
	case "items":
		items, err := read_playlist_items(db)
		if err != nil {
			err_fatal(err)
		}
		headers, display_rows := get_items_display(items)
		print_table(headers, display_rows)

	case "reset":
		clear_vid_records(db)
	case "drop":
		drop_config_table(db)
	case "insert":
		init_quota_row(db)

	case "refresh":
		refresh_token(db)

	case "quota":
		quota, err := read_quota(db)
		if err != nil {
			err_fatal(err)
		}
		msg := fmt.Sprintf("units remaining => %v", quota.quota)
		info_msg(msg)
	case "token":
		access_token := os.Getenv("ACCESS_TOKEN")
		show_gcloud_tokens(access_token)

	default:
		err = fmt.Errorf("Invalid subcommand. To see available commands, run 'scout help'")
		err_fatal(err)
	}
}
