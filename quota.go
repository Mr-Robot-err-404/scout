package main

func init_quota_map() map[string]int {
	quota_map := map[string]int{
		"get":    1,
		"search": 50,
		"insert": 100,
	}
	return quota_map
}

func read_quota() (Quota, error) {
	var quota Quota
	query := "SELECT * FROM quota"
	rows, err := db.Query(query)
	if err != nil {
		return quota, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&quota.id, &quota.quota, &quota.timestamp, &quota.last_refresh)
		if err != nil {
			return quota, err
		}
	}
	return quota, nil
}

func update_quota(units *int) {
	err := queries.Update_quota(ctx, int64(*units))
	if err != nil {
		err_fatal(err)
	}
}

func init_quota_row() error {
	err := queries.Init_quota_row(ctx)
	if err != nil {
		return err
	}
	return nil
}

func drop_quota_table() {
	query := "DROP TABLE quota"
	_, err := db.Exec(query)
	if err != nil {
		err_fatal(err)
	}
	success_msg("dropped table quota")
}
