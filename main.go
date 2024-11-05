package main

import (
	"flag"
	"fmt"
	"os"
	"scout/scout_db"

	"github.com/joho/godotenv"
)

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
	err = connect_db("./app.db")
	if err != nil {
		err_fatal(err)
	}
	cmd := flag.NewFlagSet("create_cmd", flag.ExitOnError)
	create_flag := cmd.String("add", "", "add")
	delete_flag := cmd.String("delete", "", "delete")
	edit_flag := cmd.String("edit", "", "edit")

	config_cmd := flag.NewFlagSet("config_cmd", flag.ExitOnError)
	format_flag := config_cmd.String("format", "", "format")
	category_flag := config_cmd.String("category", "", "category")
	max_flag := config_cmd.String("max", "", "max")

	// TODO:
	// compare states: exists, items.
	// if !exists: remove local playlist
	// if items != config.items, add to queue

	switch os.Args[1] {
	case "cli":
		str := parse_html_str("Chess Opening Traps | Philidor, King&#39;s Indian, London, Caro-Kann | GM Naroditsky&#39;s DYI Speedrun")
		fmt.Println(str)

	case "setup":
		err := create_tables()
		if err != nil {
			err_fatal(err)
		}
		success_msg("created tables")
		err = init_quota_row()
		if err != nil {
			err_fatal(err)
		}
		success_msg("quota row initialized")
	case "cron":
		check_token()
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		quota, err := read_quota()
		if err != nil {
			err_fatal(err)
		}
		units := quota.quota
		defer update_quota(&units)

		updated, err := cron_job(api_key, access_token, &units)
		if err != nil {
			err_fatal(err)
		}
		update_vids(updated)
		update_items(updated)
		success_resp()

	case "video", "vid":
		videos, err := read_videos()
		if err != nil {
			err_fatal(err)
		}
		headers, display_rows := get_video_display(videos)
		print_table(headers, display_rows)

	case "channel", "chan":
		if len(os.Args) == 2 {
			channels := read_channels()
			headers, display_rows := get_channel_display(channels)
			print_table(headers, display_rows)
			return
		}
		cmd.Parse(os.Args[2:])
		if len(*delete_flag) != 0 {
			tag, err, _ := find_channel(*delete_flag)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					err_msg("no channel found with that tag")
					return
				}
				err_fatal(err)
			}
			log := "delete channel"
			err = queries.Delete_channel_row(ctx, tag)
			if err != nil {
				err_msg(log)
				err_fatal(err)
			}
			success_msg(log)
			return
		}
		if len(*edit_flag) != 0 {
			channel_id := *edit_flag
			err := edit_channel(channel_id)
			if err != nil {
				err_fatal(err)
			}
			return
		}
		if len(*create_flag) == 0 {
			err_msg("invalid flags for channel cmd")
			return
		}
		_, _, exists := find_channel(*create_flag)
		if exists {
			info_msg_fatal("channel is already tracked")
		}
		config, err := read_config_file()
		if err != nil {
			err_fatal(err)
		}
		log := "add channel"
		load(log)
		quota, err := read_quota()

		if err != nil {
			err_fatal(err)
		}
		units := quota.quota
		defer update_quota(&units)

		key := os.Getenv("API_KEY")
		item, err := get_channel_ID(*create_flag, key, &units)
		if err != nil {
			err_msg(log)
			err_fatal(err)
		}
		id, title, real_tag := item[0], item[1], item[2]
		params := scout_db.Create_channel_row_params{ChannelID: id, Name: title, Tag: real_tag, Category: config.category}
		err = queries.Create_channel_row(ctx, params)

		if err != nil {
			err_msg(log)
			err_fatal(err)
		}
		success_msg(log)

	case "playlist", "play":
		if len(os.Args) == 2 {
			playlists := read_playlists()
			headers, display_rows := get_playlist_display(playlists)
			print_table(headers, display_rows)
			return
		}
		cmd.Parse(os.Args[2:])

		if len(*delete_flag) != 0 {
			playlist_id := *delete_flag

			_, err := queries.Delete_playlist(ctx, playlist_id)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					err_msg("no playlists match that ID")
					return
				}
				err_fatal(err)
			}
			log := "delete playlist"
			load(log)
			check_token()

			api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
			err = delete_remote_playlist(playlist_id, api_key, access_token)
			if err != nil {
				err_msg(log)
				err_fatal(err)
			}
			success_msg(log)
			return
		}
		if len(*edit_flag) != 0 {
			playlist_id := *edit_flag
			err := edit_playlist(playlist_id)
			if err != nil {
				err_fatal(err)
			}
			return
		}
		if len(*create_flag) == 0 {
			err := fmt.Errorf("no valid flags given for play cmd")
			err_fatal(err)
		}
		query := get_user_input("Enter search terms: ", true)
		filter := get_user_input("Filter: ", false)

		quota, err := read_quota()
		if err != nil {
			err_fatal(err)
		}
		units := quota.quota
		defer update_quota(&units)

		check_token()
		config, err := read_config_file()
		if err != nil {
			err_fatal(err)
		}
		q := csv_string(query)
		f := csv_string(filter)

		items, err := select_playlist_items(query, filter, &units, config.max_items, config.format, config.category)
		if err != nil {
			err_fatal(err)
		}
		if len(items) == 0 {
			info_msg_fatal("no matching items found")
		}
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		resp := create_playlist(*create_flag, api_key, access_token, &units)
		videos, c := populate_playlist(resp.Id, &units, items, api_key, access_token)

		params := scout_db.Add_playlist_row_params{PlaylistID: resp.Id, Name: resp.Snippet.Title, Q: q, Filter: f, Items: int64(c), Category: config.category, Format: config.format}
		queries.Add_playlist_row(ctx, params)

		add_vid_rows(videos)

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

	case "reset":
		clear_vid_records()
	case "drop":
		drop_config_table()
	case "insert":
		err := init_quota_row()
		if err != nil {
			err_fatal(err)
		}
		success_msg("quota table initialized")

	case "refresh":
		err := refresh_token()
		if err != nil {
			err_fatal(err)
		}
		success_msg("refresh token")

	case "quota":
		quota, err := read_quota()
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
