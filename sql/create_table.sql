CREATE TABLE IF NOT EXISTS channel (
        id SERIAL,
	channel_id VARCHAR(100),
	tag VARCHAR(100),
	name VARCHAR(100),
	category VARCHAR(100),
	PRIMARY KEY (id)
)
