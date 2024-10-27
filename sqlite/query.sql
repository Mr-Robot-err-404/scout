-- name: Insert_config_row :exec
INSERT INTO config (format, category, max_items) 
VALUES (?, ?, ?);

-- name: Create_channel_row :exec
INSERT INTO channel (channel_id, tag, name, category) 
VALUES (?, ?, ?, ?);

-- name: Add_playlist_row :exec
INSERT INTO playlist (playlist_id, name, q, filter, format, items, category) 
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: Insert_vid_row :exec
INSERT INTO video (video_id, title) 
VALUES (?, ?);

-- name: Delete_playlist :exec
DELETE FROM playlist 
WHERE name = ?;

-- name: Delete_channel_row :exec
DELETE FROM channel 
WHERE tag = ?;

-- name: Find_playlist_name :one
SELECT name 
FROM playlist 
WHERE name = ?;

-- name: Find_playlist_id :one
SELECT * 
FROM playlist 
WHERE playlist_id = ?;

-- name: Update_playlist :exec
UPDATE playlist
SET q = ?, filter = ?, category = ?, format = ?
WHERE playlist_id = ?;

-- name: Find_channel_row :one
SELECT tag 
FROM channel
WHERE tag = ?;

-- name: Init_quota_row :exec
INSERT INTO quota (id, quota, quota_reset, last_refresh) 
VALUES (69, 10000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: Update_quota :exec
UPDATE quota 
SET quota = ?;

-- name: Update_last_refresh :exec
UPDATE quota 
SET last_refresh = CURRENT_TIMESTAMP;
