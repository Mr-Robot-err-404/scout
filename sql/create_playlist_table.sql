CREATE TABLE IF NOT EXISTS playlist (
        id SERIAL,
	playlist_id VARCHAR(100),
	name VARCHAR(100),
	q VARCHAR(500),
	filter VARCHAR(500),
	long_format BOOLEAN DEFAULT TRUE,
	PRIMARY KEY (id)
)
