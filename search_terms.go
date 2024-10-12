package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func get_user_input(msg string, required bool) []string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		err_fatal(err)
	}
	str := scanner.Text()
	if len(str) == 0 && required {
		err := fmt.Errorf("that field is required good sir")
		err_fatal(err)
	}
	q := []string{}

	for _, s := range strings.Split(str, ",") {
		q = append(q, strings.TrimSpace(s))
	}
	return q
}
