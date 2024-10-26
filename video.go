package main

import (
	"scout/scout_db"
)

type Video struct {
	video_id string
	title    string
}

func read_videos() ([]Video, error) {
	videos := []Video{}
	query := "SELECT * FROM video;"
	rows, err := db.Query(query)
	if err != nil {
		return videos, err
	}
	defer rows.Close()
	for rows.Next() {
		var vid Video
		err = rows.Scan(&vid.video_id, &vid.title)
		if err != nil {
			return videos, err
		}
		videos = append(videos, vid)
	}
	return videos, nil
}

func add_vid_rows(videos []Video) {
	for i := range videos {
		curr := videos[i]
		err := queries.Insert_vid_row(ctx, scout_db.Insert_vid_row_params{VideoID: curr.video_id, Title: curr.title})
		if err != nil {
			err_resp(err)
		}
	}
	success_resp()
}

func drop_vid_table() {
	query := "DROP TABLE video"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table")

}

func clear_vid_records() {
	query := "DELETE FROM video"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("cleared table")
}
