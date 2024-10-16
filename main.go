package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var sp = create_spinner()

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

	// TODO: track videos & playlist items

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
			channels := readChannels(db)
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
		key := os.Getenv("API_KEY")
		log := "add channel"
		load(log)
		item, err := get_channel_ID(*create_flag, key)

		if err != nil {
			err_msg(log)
			err_fatal(err)
		}
		id, title, real_tag := item[0], item[1], item[2]
		err = createChannelRow(db, id, real_tag, title)

		if err != nil {
			err_msg(log)
			err_fatal(err)
		}
		success_msg(log + " => " + title)

	case "playlist":
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
		query := get_user_input("Enter search terms: ", true)
		filter := get_user_input("Filter: ", false)

		check_token(db)
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		q := csv_string(query)
		f := csv_string(filter)

		playlist_resp := create_playlist(db, *create_flag, q, f, api_key, access_token)
		videos := populate_playlist(db, query, filter, playlist_resp.Id)
		add_vid_rows(db, videos)

	case "table":
		err := createTable(db, "./sql/daily_quota.sql")
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
		clear_item_table(db)
	case "drop":
		drop_quota_table(db)
	case "insert":
		init_quota_row(db)

	case "refresh":
		refresh_token(db)

	case "quota":
		quota, err := read_quota(db)
		if err != nil {
			err_fatal(err)
		}
		fmt.Println(time.Now().Unix() - quota.last_refresh.Unix())
	case "token":
		access_token := os.Getenv("ACCESS_TOKEN")
		credentials := readCredentialsFile("../.config/gcloud/application_default_credentials.json")
		fmt.Println("----------------------------------------------")
		fmt.Printf("REFRESH_TOKEN %v\n", credentials.Refresh_token)
		fmt.Println("----------------------------------------------")
		fmt.Printf("CLIENT_ID     %v\n", credentials.Client_id)
		fmt.Println("----------------------------------------------")
		fmt.Printf("CLIENT_SECRET %v\n", credentials.Client_secret)
		fmt.Println("----------------------------------------------")
		fmt.Printf("ACCESS_TOKEN %v\n", access_token)
		fmt.Println("----------------------------------------------")

	default:
		err = fmt.Errorf("Invalid subcommand. To see available commands, run 'scout help'")
		err_fatal(err)
	}
}
