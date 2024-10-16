CREATE TABLE IF NOT EXISTS playlist_item (
        id SERIAL,
	video_id VARCHAR(100),
	playlist_id VARCHAR(100),
	channel_id VARCHAR(100),
	PRIMARY KEY (id),
	FOREIGN KEY (video_id) REFERENCES video(video_id),
	FOREIGN KEY (playlist_id) REFERENCES playlist(playlist_id), 
	FOREIGN KEY (channel_id) REFERENCES channel(channel_id)
)
