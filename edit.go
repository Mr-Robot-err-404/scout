package main

func edit_playlist(playlist_id string) error {
	playlist := Playlist{}
	query := `SELECT * FROM playlist 
	WHERE playlist_id = $1`

	row := db.QueryRow(query, playlist_id)
	err := row.Scan(&playlist.playlist_id, &playlist.name, &playlist.q, &playlist.filter, &playlist.format, &playlist.items, &playlist.category)
	if err != nil {
		return err
	}
	print_title_with_bg("edit " + playlist.name)
	str, err := edit_user_input("Search query: ", playlist.q)

	if err != nil {
		err_fatal(err)
	}
	original := parse_input(playlist.q)
	q := parse_input(str)

	if is_edited(original, q) {
		// TODO: write query to update field
	}
	return nil
}

func is_edited(original []string, input []string) bool {
	if len(original) != len(input) {
		return true
	}
	for i := range original {
		if original[i] != input[i] {
			return true
		}
	}
	return false
}
