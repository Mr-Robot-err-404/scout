package main

import (
	"database/sql"
	"fmt"
)

func scrape_channels(db *sql.DB) {
	channels := readChannels(db)
	fmt.Println(channels)
}
