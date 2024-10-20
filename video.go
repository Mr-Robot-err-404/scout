package main

import "database/sql"

type Video struct {
	video_id string
	title    string
}

func read_videos(db *sql.DB) ([]Video, error) {
	videos := []Video{}
	query := readSQLFile("./sql/read_all_videos.sql")
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

func insert_vid_row(db *sql.DB, video_id string, title string) error {
	query := readSQLFile("./sql/create_video.sql")
	_, err := db.Exec(query, video_id, title)
	if err != nil {
		return err
	}
	return nil
}

func add_vid_rows(db *sql.DB, videos []Video) {
	for i := range videos {
		curr := videos[i]
		err := insert_vid_row(db, curr.video_id, curr.title)
		if err != nil {
			err_resp(err)
		}
	}
	success_resp()
}

func drop_vid_table(db *sql.DB) {
	query := "DROP TABLE video"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table")

}

func clear_vid_records(db *sql.DB) {
	query := "DELETE FROM video"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("cleared table")
}
