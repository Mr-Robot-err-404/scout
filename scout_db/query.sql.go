// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package scout_db

import (
	"context"
	"time"
)

const add_playlist_row = `-- name: Add_playlist_row :exec
INSERT INTO playlist (playlist_id, name, q, filter, format, items, category) 
VALUES (?, ?, ?, ?, ?, ?, ?)
`

type Add_playlist_row_params struct {
	PlaylistID string
	Name       string
	Q          string
	Filter     string
	Format     string
	Items      int64
	Category   string
}

func (q *Queries) Add_playlist_row(ctx context.Context, arg Add_playlist_row_params) error {
	_, err := q.db.ExecContext(ctx, add_playlist_row,
		arg.PlaylistID,
		arg.Name,
		arg.Q,
		arg.Filter,
		arg.Format,
		arg.Items,
		arg.Category,
	)
	return err
}

const channels_by_category = `-- name: Channels_by_category :many
SELECT channel_id, tag, name, category
FROM channel 
WHERE category = ?
`

func (q *Queries) Channels_by_category(ctx context.Context, category string) ([]Channel, error) {
	rows, err := q.db.QueryContext(ctx, channels_by_category, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Channel
	for rows.Next() {
		var i Channel
		if err := rows.Scan(
			&i.ChannelID,
			&i.Tag,
			&i.Name,
			&i.Category,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const create_channel_row = `-- name: Create_channel_row :exec
INSERT INTO channel (channel_id, tag, name, category) 
VALUES (?, ?, ?, ?)
`

type Create_channel_row_params struct {
	ChannelID string
	Tag       string
	Name      string
	Category  string
}

func (q *Queries) Create_channel_row(ctx context.Context, arg Create_channel_row_params) error {
	_, err := q.db.ExecContext(ctx, create_channel_row,
		arg.ChannelID,
		arg.Tag,
		arg.Name,
		arg.Category,
	)
	return err
}

const delete_channel_row = `-- name: Delete_channel_row :exec
DELETE FROM channel 
WHERE tag = ?
`

func (q *Queries) Delete_channel_row(ctx context.Context, tag string) error {
	_, err := q.db.ExecContext(ctx, delete_channel_row, tag)
	return err
}

const delete_playlist = `-- name: Delete_playlist :one
DELETE FROM playlist 
WHERE playlist_id = ?
RETURNING playlist_id, name, q, "filter", format, items, category
`

func (q *Queries) Delete_playlist(ctx context.Context, playlistID string) (Playlist, error) {
	row := q.db.QueryRowContext(ctx, delete_playlist, playlistID)
	var i Playlist
	err := row.Scan(
		&i.PlaylistID,
		&i.Name,
		&i.Q,
		&i.Filter,
		&i.Format,
		&i.Items,
		&i.Category,
	)
	return i, err
}

const find_channel = `-- name: Find_channel :one
SELECT channel_id, tag, name, category 
FROM channel
WHERE channel_id = ?
`

func (q *Queries) Find_channel(ctx context.Context, channelID string) (Channel, error) {
	row := q.db.QueryRowContext(ctx, find_channel, channelID)
	var i Channel
	err := row.Scan(
		&i.ChannelID,
		&i.Tag,
		&i.Name,
		&i.Category,
	)
	return i, err
}

const find_channel_tag = `-- name: Find_channel_tag :one
SELECT tag
FROM channel
WHERE tag = ?
`

func (q *Queries) Find_channel_tag(ctx context.Context, tag string) (string, error) {
	row := q.db.QueryRowContext(ctx, find_channel_tag, tag)
	err := row.Scan(&tag)
	return tag, err
}

const find_playlist = `-- name: Find_playlist :one
SELECT playlist_id, name, q, "filter", format, items, category 
FROM playlist 
WHERE playlist_id = ?
`

func (q *Queries) Find_playlist(ctx context.Context, playlistID string) (Playlist, error) {
	row := q.db.QueryRowContext(ctx, find_playlist, playlistID)
	var i Playlist
	err := row.Scan(
		&i.PlaylistID,
		&i.Name,
		&i.Q,
		&i.Filter,
		&i.Format,
		&i.Items,
		&i.Category,
	)
	return i, err
}

const find_playlist_name = `-- name: Find_playlist_name :one
SELECT name 
FROM playlist 
WHERE name = ?
`

func (q *Queries) Find_playlist_name(ctx context.Context, name string) (string, error) {
	row := q.db.QueryRowContext(ctx, find_playlist_name, name)
	err := row.Scan(&name)
	return name, err
}

const init_quota_row = `-- name: Init_quota_row :exec
INSERT INTO quota (id, quota, quota_reset, last_refresh) 
VALUES (69, 10000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
`

func (q *Queries) Init_quota_row(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, init_quota_row)
	return err
}

const insert_config_row = `-- name: Insert_config_row :exec
INSERT INTO config (format, category, max_items) 
VALUES (?, ?, ?)
`

type Insert_config_row_params struct {
	Format   string
	Category string
	MaxItems int64
}

func (q *Queries) Insert_config_row(ctx context.Context, arg Insert_config_row_params) error {
	_, err := q.db.ExecContext(ctx, insert_config_row, arg.Format, arg.Category, arg.MaxItems)
	return err
}

const insert_vid_row = `-- name: Insert_vid_row :exec
INSERT INTO video (video_id, title) 
VALUES (?, ?)
`

type Insert_vid_row_params struct {
	VideoID string
	Title   string
}

func (q *Queries) Insert_vid_row(ctx context.Context, arg Insert_vid_row_params) error {
	_, err := q.db.ExecContext(ctx, insert_vid_row, arg.VideoID, arg.Title)
	return err
}

const update_channel_category = `-- name: Update_channel_category :exec
UPDATE channel
SET category = ?
WHERE channel_id = ?
`

type Update_channel_category_params struct {
	Category  string
	ChannelID string
}

func (q *Queries) Update_channel_category(ctx context.Context, arg Update_channel_category_params) error {
	_, err := q.db.ExecContext(ctx, update_channel_category, arg.Category, arg.ChannelID)
	return err
}

const update_last_refresh = `-- name: Update_last_refresh :exec
UPDATE quota 
SET last_refresh = CURRENT_TIMESTAMP
`

func (q *Queries) Update_last_refresh(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, update_last_refresh)
	return err
}

const update_playlist = `-- name: Update_playlist :exec
UPDATE playlist
SET q = ?, filter = ?, category = ?, format = ?
WHERE playlist_id = ?
`

type Update_playlist_params struct {
	Q          string
	Filter     string
	Category   string
	Format     string
	PlaylistID string
}

func (q *Queries) Update_playlist(ctx context.Context, arg Update_playlist_params) error {
	_, err := q.db.ExecContext(ctx, update_playlist,
		arg.Q,
		arg.Filter,
		arg.Category,
		arg.Format,
		arg.PlaylistID,
	)
	return err
}

const update_playlist_item_count = `-- name: Update_playlist_item_count :exec
UPDATE playlist
SET items = ?
WHERE playlist_id = ?
`

type Update_playlist_item_count_params struct {
	Items      int64
	PlaylistID string
}

func (q *Queries) Update_playlist_item_count(ctx context.Context, arg Update_playlist_item_count_params) error {
	_, err := q.db.ExecContext(ctx, update_playlist_item_count, arg.Items, arg.PlaylistID)
	return err
}

const update_quota = `-- name: Update_quota :exec
UPDATE quota 
SET quota = ?, quota_reset = ?
`

type Update_quota_params struct {
	Quota      int64
	QuotaReset time.Time
}

func (q *Queries) Update_quota(ctx context.Context, arg Update_quota_params) error {
	_, err := q.db.ExecContext(ctx, update_quota, arg.Quota, arg.QuotaReset)
	return err
}
