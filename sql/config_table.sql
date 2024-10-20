CREATE TABLE IF NOT EXISTS config (
	id SERIAL,
	format VARCHAR(100),
	category VARCHAR(100),
	max_items INTEGER,
	PRIMARY KEY (id)
) 
