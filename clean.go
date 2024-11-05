package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func clean_data() {
	var s SearchResp
	data, err := os.ReadFile("./samples/rosen.json")
	if err != nil {
		err_fatal(err)
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		err_fatal(err)
	}
	curr := s.Items[0]
	fmt.Println(curr)
}

func clean_resp(s string) string {
	const key = "&#39;"
	slice := strings.Split(s, key)
	if len(slice) == 0 {
		return ""
	}
	str := slice[0]
	for i := 1; i < len(slice); i++ {
		str += "'" + slice[i]
	}
	return str
}
