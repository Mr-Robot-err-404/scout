package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Message struct {
	Name string
	Body string
	Time int64
}

var sp = create_spinner()

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		help_txt, err := os.ReadFile("./help.txt")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(help_txt))
		return
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	db := setup_db()
	defer db.Close()

	cmd := flag.NewFlagSet("create_cmd", flag.ExitOnError)
	create_flag := cmd.String("add", "", "add")
	delete_flag := cmd.String("delete", "", "delete")

	// TODO: replace all logs with a state msg (load, success, err, info)

	switch os.Args[1] {
	case "cli":
		logging_time()

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
			} else {
				log := "delete channel"
				load(log)
				err = deleteRow(db, tag)
				if err != nil {
					err_msg(log)
					err_resp(err)
					return
				}
				success_msg(log)
			}
			return
		}
		if len(*create_flag) == 0 {
			err_msg("no channel tag provided")
			return
		}
		_, exists := find_row(db, *create_flag, "./sql/read_row.sql")
		if exists {
			info_msg("Channel is already tracked")
			return
		}
		key := os.Getenv("API_KEY")
		item, err := get_channel_ID(*create_flag, key)

		if err != nil {
			err_fatal(err)
		}
		id, title, real_tag := item[0], item[1], item[2]
		err = createChannelRow(db, id, real_tag, title)
		if err != nil {
			err_fatal(err)
		}

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
				os.Exit(1)
			}
			os.Exit(0)
		}
		if len(*create_flag) == 0 {
			log.Fatal("playlist name not provided")
		}
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		query := get_user_input("Enter search terms: ")
		filter := get_user_input("Filter: ")

		q := csv_string(query)
		f := csv_string(filter)

		playlist_resp := create_playlist(db, *create_flag, q, f, api_key, access_token)
		populate_playlist(db, query, filter, playlist_resp.Id)

	case "create_table":
		createTable(db, "./sql/create_playlist_table.sql")
	case "delete_table":
		deleteTable(db, "./sql/delete_playlist_table.sql")
	case "refresh":
		refresh_quota(db)
	case "quota":
		quota := read_quota(db)
		fmt.Println(time.Now().Unix() - quota.timestamp.Unix())
	case "token":
		credentials := readCredentialsFile("../.config/gcloud/application_default_credentials.json")
		fmt.Println("----------------------------------------------")
		fmt.Printf("REFRESH TOKEN %v\n", credentials.Refresh_token)
		fmt.Println("----------------------------------------------")
		fmt.Printf("CLIENT_ID %v\n", credentials.Client_id)
		fmt.Println("----------------------------------------------")
		fmt.Printf("CLIENT_SECRET %v\n", credentials.Client_secret)
		fmt.Println("----------------------------------------------")
	case "insert":
		insert_row(db)
	default:
		log.Fatal("Invalid subcommand. To see usable commands, use 'cli help'")
	}
}
