package main

import (
	"context"
	"database/sql"
	_ "embed"
	"scout/scout_db"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var db *sql.DB
var queries *scout_db.Queries
var ctx context.Context

//go:embed sqlite/schema.sql
var ddl string

func connect_db(db_path string) error {
	var err error
	db, err = sql.Open("sqlite", db_path)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	ctx = context.Background()
	queries = scout_db.New(db)
	return nil
}

func create_tables() error {
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}
	return nil
}
