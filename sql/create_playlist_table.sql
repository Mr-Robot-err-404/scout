CREATE TABLE IF NOT EXISTS playlist (
	playlist_id VARCHAR(100),
	name VARCHAR(100),
	q VARCHAR(500),
	filter VARCHAR(500),
	long_format BOOLEAN DEFAULT TRUE,
	PRIMARY KEY (playlist_id)
)
