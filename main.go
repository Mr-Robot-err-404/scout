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
	populate := create_cmd.Bool("populate", false, "populate")

	// TODO: ask a second prompt if first condition is valid

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
		key := os.Getenv("SEARCH_KEY_1")
		item, err := getChannelID(tag, key)

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
	case "scan":
		q := get_search_terms()
		fmt.Println(q)

	case "playlist":
		if len(os.Args) == 2 {
			playlists := read_playlists(db)
			fmt.Println(playlists)
			return
		}
		insert_key, access_token := os.Getenv("PLAYLIST_KEY_1"), os.Getenv("ACCESS_TOKEN")
		create_cmd.Parse(os.Args[2:])

		if *populate {
			scrape_channels(db)
			os.Exit(0)
		}

		if len(*playlist_name) == 0 {
			log.Fatal("playlist name not provided")
		}

		item := create_playlist(db, *playlist_name, insert_key, access_token)
		fmt.Printf("created playlist: %v", item)

	case "create_table":
		createTable(db, "./sql/create_playlist_table.sql")
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
