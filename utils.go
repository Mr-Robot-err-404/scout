package main

import (
	"fmt"
	"html"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func parse_html_str(s string) string {
	return html.UnescapeString(s)
}

func find_channel(search_term string) (string, error, bool) {
	s := search_term[:]

	if string(s[0]) != "@" {
		s = "@" + s
	}
	tag, err := queries.Find_channel_tag(ctx, s)
	if err != nil {
		return "", err, false
	}
	return tag, nil, true
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

func extract_pt_time(ts time.Time) time.Time {
	time_diff := 8 * 3600
	pt_time := ts.Unix() - int64(time_diff)
	return time.Unix(pt_time, 0)
}

func parse_query(q string) string {
	return url.PathEscape(q)
}

func parse_input(str string) []string {
	q := []string{}
	for _, s := range strings.Split(str, ",") {
		q = append(q, strings.TrimSpace(s))
	}
	return q
}

func convert_and_parse(q []string) string {
	if len(q) == 0 {
		return ""
	}
	query := q[0]
	if len(q) == 1 {
		return parse_query(query)
	}
	for i := 1; i < len(q); i++ {
		s := q[i]
		query += "|"
		query += s
	}
	return parse_query(query)
}

func validate_config_flags(format string, category string, max_items string, track string) (map[string]string, error) {
	config_options := make(map[string]string)
	if len(format) > 0 {
		if format == "short" || format == "medium" || format == "long" {
			config_options["format"] = format
		} else {
			err := fmt.Errorf("invalid config option: 'format'. Accepted values: short || medium || long")
			return config_options, err
		}
	}
	if len(category) > 0 {
		config_options["category"] = category
	}
	if len(max_items) > 0 {
		num, err := strconv.Atoi(max_items)
		if err != nil {
			return config_options, err
		}
		if num < 0 || num > 50 {
			err := fmt.Errorf("invalid config option: 'max'. Accepted values: 0 -> 50")
			return config_options, err
		}
		config_options["max"] = max_items
	}
	if len(track) > 0 {
		if track != "on" && track != "off" {
			err := fmt.Errorf("invalid config option: 'track'. Accepted values: on || off")
			return config_options, err
		}
		config_options["track"] = track
	}
	return config_options, nil
}

func show_gcloud_tokens(access_token string) {
	credentials := read_credentials_file("../.config/gcloud/application_default_credentials.json")
	fmt.Println("----------------------------------------------")
	fmt.Printf("REFRESH_TOKEN %v\n", credentials.Refresh_token)
	fmt.Println("----------------------------------------------")
	fmt.Printf("CLIENT_ID     %v\n", credentials.Client_id)
	fmt.Println("----------------------------------------------")
	fmt.Printf("CLIENT_SECRET %v\n", credentials.Client_secret)
	fmt.Println("----------------------------------------------")
	fmt.Printf("ACCESS_TOKEN %v\n", access_token)
	fmt.Println("----------------------------------------------")

}

func get_channel_IDs(channels []Channel) []string {
	IDs := []string{}
	for i := range channels {
		curr := channels[i]
		IDs = append(IDs, curr.channel_id)
	}
	return IDs
}
