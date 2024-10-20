CREATE TABLE IF NOT EXISTS playlist (
	playlist_id VARCHAR(100),
	name VARCHAR(100),
	q VARCHAR(500),
	filter VARCHAR(500),
	format VARCHAR(100),
	items INTEGER,
	category VARCHAR(100),
	PRIMARY KEY (playlist_id)
)
