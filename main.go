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

func main() {
	help_txt, err := os.ReadFile("./help.txt")
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Print(string(help_txt))
		return
	}
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	db := setup_db()
	defer db.Close()

	create_cmd := flag.NewFlagSet("create_cmd", flag.ExitOnError)
	playlist_name := create_cmd.String("create", "", "create")
	delete_flag := create_cmd.String("delete", "", "delete")

	// TODO: insert a video into an existing playlist

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 3 {
			log.Fatal("no channel tag provided")
		}
		tag := os.Args[2]
		_, exists := find_row(db, tag, "./sql/read_row.sql")
		if exists {
			log.Println("Channel is already tracked ;)")
			os.Exit(0)
		}
		key := os.Getenv("API_KEY")
		item, err := get_channel_ID(tag, key)

		if err != nil {
			log.Fatal(err)
		}
		id, title, real_tag := item[0], item[1], item[2]
		createChannelRow(db, id, real_tag, title)
	case "remove":
		if len(os.Args) != 3 {
			log.Fatal("no channel tag provided")
		}
		tag := os.Args[2]
		tag, exists := find_row(db, tag, "./sql/read_row.sql")
		if !exists {
			log.Println("no channel found")
			os.Exit(0)
		}
		deleteRow(db, tag)
	case "cli":
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		res := insert_playlist_item("PL-vGMW-bu9eXQ3mWWHRY5hJ6xLVhhgRFh", "tO7CCP7liwI", api_key, access_token)
		fmt.Printf("added playlist item: %v", res.Snippet.Title)

	case "playlist":
		if len(os.Args) == 2 {
			playlists := read_playlists(db)
			fmt.Println(playlists)
			return
		}
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		create_cmd.Parse(os.Args[2:])

		if len(*delete_flag) != 0 {
			err := deletePlaylist(db, *delete_flag)
			if err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}

		if len(*playlist_name) == 0 {
			log.Fatal("playlist name not provided")
		}
		terms := get_user_search_terms()
		q := csv_string(terms)

		item := create_playlist(db, *playlist_name, q, api_key, access_token)
		fmt.Printf("created playlist: %v", item)

	case "create_table":
		createTable(db, "./sql/create_video_table.sql")
	case "delete_table":
		deleteTable(db, "./sql/delete_playlist_table.sql")
	case "refresh":
		refresh_quota(db)
	case "quota":
		quota := read_quota(db)
		fmt.Println(time.Now().Unix() - quota.timestamp.Unix())
	case "read":
		channels := readChannels(db)
		fmt.Println(channels)
	case "insert":
		insert_row(db)
	default:
		log.Fatal("Invalid subcommand. To see usable commands, use 'cli help'")
	}
}
