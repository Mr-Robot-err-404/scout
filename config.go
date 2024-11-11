package main

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	format    string
	category  string
	max_items int
	track     string
}

func init_config_table() {
	query := readSQLFile("./sql/config_table.sql")
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	config := Config{
		format:    "long",
		category:  "chess",
		max_items: 5,
	}
	err = insert_config_row(config)
	if err != nil {
		err_fatal(err)
	}
	setup_config_file(config)
	success_msg("config setup")
}

func setup_config_file(config Config) error {
	env, err := get_config_map()
	if err != nil {
		return err
	}
	env["format"] = config.format
	env["category"] = config.category
	env["max"] = strconv.Itoa(config.max_items)
	err = godotenv.Write(env, "./.config")
	if err != nil {
		return err
	}
	return nil
}

func update_config_file(config map[string]string) error {
	env, err := get_config_map()
	if err != nil {
		return err
	}
	for key, value := range config {
		env[key] = value
	}
	err = godotenv.Write(env, "./.config")
	if err != nil {
		return err
	}
	return nil
}

func get_config_map() (map[string]string, error) {
	var env_map map[string]string
	env_map, err := godotenv.Read("./.config")
	if err != nil {
		return env_map, err
	}
	return env_map, nil
}

func set_default_config() error {
	var default_map = map[string]string{
		"format":   "medium",
		"category": "chess",
		"max":      "10",
		"track":    "on",
	}
	err := godotenv.Write(default_map, "./.config")
	if err != nil {
		return err
	}
	return nil
}

func insert_config_row(config Config) error {
	query := readSQLFile("./sql/config_row.sql")
	_, err := db.Exec(query, config.format, config.category, config.max_items)
	if err != nil {
		return err
	}
	return nil
}

func print_config(config Config) {
	format_msg := fmt.Sprintf("video format | %v", config.format)
	category_msg := fmt.Sprintf("category     | %v", config.category)
	max_msg := fmt.Sprintf("max items    | %v", config.max_items)
	track_msg := fmt.Sprintf("tracking     | %v", config.track)

	info_msg(format_msg)
	info_msg(category_msg)
	info_msg(max_msg)
	info_msg(track_msg)
}

func read_config_file() (Config, error) {
	var config Config
	env, err := get_config_map()
	if err != nil {
		return config, err
	}
	max_items, err := strconv.Atoi(env["max"])
	if err != nil {
		return config, err
	}
	config.format = env["format"]
	config.category = env["category"]
	config.track = env["track"]
	config.max_items = max_items
	return config, nil
}

func read_config() (Config, error) {
	config := Config{}
	query := "SELECT * FROM config"
	rows, err := db.Query(query)
	var id int
	if err != nil {
		return config, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &config.format, &config.category, &config.max_items)
		if err != nil {
			return config, err
		}
	}
	return config, nil

}

func drop_config_table() {
	query := "DROP TABLE config"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table")

}
