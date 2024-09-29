package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

func get_user_search_terms() []string {
	scanner := bufio.NewScanner(os.Stdin)
	msg := "Enter search terms: "
	fmt.Fprint(os.Stdout, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	str := scanner.Text()
	if len(str) == 0 {
		log.Fatal("absolutely spiffing choice, sir...")
	}
	q := []string{}

	for _, s := range strings.Split(str, ",") {
		q = append(q, strings.TrimSpace(s))
	}
	return q
}
func search_remote_channels(db *sql.DB, q []string) {
	res, err := search_and_destroy("", "UCXy10-NEFGxQ3b4NVrzHw1Q")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("page info: %v\n", res.PageInfo)
	for i := 0; i < len(res.Items); i++ {
		item := res.Items[i]
		fmt.Printf("videoID : %v\n", item.Id.VideoId)
		fmt.Printf("title : %v\n", item.Snippet.Title)
		fmt.Printf("description : %v\n", item.Snippet.Description)
		fmt.Printf("channelID : %v\n", item.Snippet.ChannelId)
		fmt.Printf("channel title: %v\n", item.Snippet.ChannelTitle)
	}
}
func csv_string(q []string) string {
	csv_line := ""
	for i := range q {
		str := q[i]
		if i == 0 {
			csv_line += str
			continue
		}
		csv_line += "," + str
	}
	return csv_line
}
