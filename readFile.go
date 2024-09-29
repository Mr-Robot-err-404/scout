package main

import (
	"log"
	"os"
)

func readSQLFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(data) 
}