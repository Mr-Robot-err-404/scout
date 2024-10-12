CREATE TABLE IF NOT EXISTS playlist_item (
        id SERIAL,
	vid_id INT,
	list_id INT,
	chan_id INT,
	PRIMARY KEY (id),
	FOREIGN KEY (vid_id) REFERENCES video(id),
	FOREIGN KEY (list_id) REFERENCES playlist(id), 
	FOREIGN KEY (chan_id) REFERENCES channel(id)
)
