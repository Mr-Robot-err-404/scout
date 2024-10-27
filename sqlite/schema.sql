CREATE TABLE IF NOT EXISTS config (
	id INTEGER NOT NULL,
	format VARCHAR(100) NOT NULL,
	category VARCHAR(100) NOT NULL,
	max_items INTEGER NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS playlist (
	playlist_id VARCHAR(100) NOT NULL,
	name VARCHAR(100) NOT NULL,
	q VARCHAR(500) NOT NULL,
	filter VARCHAR(500) NOT NULL,
	format VARCHAR(100) NOT NULL,
	items INTEGER NOT NULL,
	category VARCHAR(100) NOT NULL,
	PRIMARY KEY (playlist_id)
);

CREATE TABLE IF NOT EXISTS channel (
	channel_id VARCHAR(100) NOT NULL,
	tag VARCHAR(100) NOT NULL,
	name VARCHAR(100) NOT NULL,
	category VARCHAR(100) NOT NULL,
	PRIMARY KEY (channel_id)
);

CREATE TABLE IF NOT EXISTS video (
	video_id VARCHAR(100) NOT NULL,
	title VARCHAR(100) NOT NULL,
	PRIMARY KEY (video_id)
);

CREATE TABLE IF NOT EXISTS quota (
        id INTEGER NOT NULL,
	quota INTEGER NOT NULL,
	quota_reset TIMESTAMP NOT NULL,
	last_refresh TIMESTAMP NOT NULL,
	PRIMARY KEY (id)
);