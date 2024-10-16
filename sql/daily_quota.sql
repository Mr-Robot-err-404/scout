CREATE TABLE IF NOT EXISTS quota (
        id SERIAL,
	quota INTEGER, 
	quota_reset TIMESTAMP,
	last_refresh TIMESTAMP,
	PRIMARY KEY (id)
)
