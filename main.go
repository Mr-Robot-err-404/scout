package main

import (
	"flag"
	"fmt"
	"os"
	"scout/scout_db"
	"time"

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
	detach_flag := cmd.String("detach", "", "detach")
	attach_flag := cmd.String("attach", "", "attach")

	config_cmd := flag.NewFlagSet("config_cmd", flag.ExitOnError)
	format_flag := config_cmd.String("format", "", "format")
	category_flag := config_cmd.String("category", "", "category")
	max_flag := config_cmd.String("max", "", "max")
	track_flag := config_cmd.String("track", "", "track")
	default_flag := config_cmd.Bool("default", false, "default")

	switch os.Args[1] {
	case "cli":
		ts := time.Now()
		fmt.Println(is_quota_reset(extract_pt_time(ts)))

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
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		token := check_token()
		if len(token) > 0 {
			access_token = token
		}
		quota, err := read_quota()
		if err != nil {
			err_fatal(err)
		}
		config, err := read_config_file()
		if err != nil {
			err_fatal(err)
		}
		units := quota.quota
		ts := quota.quota_reset
		if is_quota_reset(quota.quota_reset) {
			units = 10000
			ts = time.Now()
		}
		defer update_quota(&units, ts)

		updated, err := cron_job(api_key, access_token, &units, config)
		if err != nil {
			err_fatal(err)
		}
		if config.track == "on" {
			update_vids(updated)
		}
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
		ts := quota.quota_reset

		if is_quota_reset(quota.quota_reset) {
			units = 10000
			ts = time.Now()
		}
		defer update_quota(&units, ts)

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

		if len(*detach_flag) > 0 {
			playlist_id := *delete_flag
			_, err := queries.Delete_playlist(ctx, playlist_id)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					info_msg("no playlists match that ID")
					return
				}
				err_fatal(err)
			}
			success_msg("detach playlist")
			return
		}
		if len(*attach_flag) > 0 {
			// TODO: take remote playlist ID
			// check if id is already tracked
			// fetch items -> []video_id
			// compare with, filter & update vid table
			// create local playlist and update details
			// no need to populate
		}

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

			api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
			token := check_token()
			if len(token) > 0 {
				access_token = token
			}
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
		ts := quota.quota_reset
		if is_quota_reset(quota.quota_reset) {
			units = 10000
			ts = time.Now()
		}
		defer update_quota(&units, ts)

		token := check_token()
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		if len(token) > 0 {
			access_token = token
		}
		config, err := read_config_file()
		if err != nil {
			err_fatal(err)
		}
		q := csv_string(query)
		f := csv_string(filter)

		log := "scrape channels"
		load(log)
		items, err_queue := select_playlist_items(query, filter, &units, config.max_items, config.format, config.category)
		if len(err_queue) != 0 {
			err_msg(log)
			log_err_queue(err_queue)
			return
		}
		success_msg(log)
		if len(items) == 0 {
			info_msg_fatal("no matching items found")
		}
		resp := create_playlist(*create_flag, api_key, access_token, &units)
		log = "populate playlist"
		load(log)
		videos, c, err_queue := populate_playlist(resp.Id, &units, items, api_key, access_token)

		if len(items) == len(err_queue) {
			err_msg(log)
			log_err_queue(err_queue)
			return
		}
		success_msg(log)
		log_err_queue(err_queue)

		params := scout_db.Add_playlist_row_params{PlaylistID: resp.Id, Name: resp.Snippet.Title, Q: q, Filter: f, Items: int64(c), Category: config.category, Format: config.format}
		queries.Add_playlist_row(ctx, params)

		if config.track == "on" {
			add_vid_rows(videos)
		}

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

		if *default_flag {
			err := set_default_config()
			if err != nil {
				err_fatal(err)
			}
			success_msg("default config set")
			return
		}
		config_options, err := validate_config_flags(*format_flag, *category_flag, *max_flag, *track_flag)
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
		_, err := refresh_token()
		if err != nil {
			err_fatal(err)
		}
		success_msg("refresh token")

	case "quota":
		quota, err := read_quota()
		if err != nil {
			err_fatal(err)
		}
		units := quota.quota

		if is_quota_reset(quota.quota_reset) {
			units = 10000
			defer update_quota(&units, time.Now())
		}
		msg := fmt.Sprintf("quota units | %v", units)
		info_msg(msg)
	case "token":
		access_token := os.Getenv("ACCESS_TOKEN")
		show_gcloud_tokens(access_token)

	default:
		err = fmt.Errorf("Invalid subcommand. To see available commands, run 'scout help'")
		err_fatal(err)
	}
}
