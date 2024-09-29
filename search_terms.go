package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func get_search_terms() []string {
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
