package main

import (
	"scout/scout_db"
	"time"
)

var quota_map = init_quota_map()

func init_quota_map() map[string]int {
	quota_map := map[string]int{
		"get":    1,
		"search": 50,
		"insert": 100,
	}
	return quota_map
}

func is_quota_reset(ts time.Time) bool {
	prev_y, prev_m, prev_d := extract_pt_time(ts).Date()
	curr_y, curr_m, curr_d := extract_pt_time(time.Now()).Date()
	if prev_y != curr_y || prev_m != curr_m || prev_d != curr_d {
		return true
	}
	return false
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
		err = rows.Scan(&quota.id, &quota.quota, &quota.quota_reset, &quota.last_refresh)
		if err != nil {
			return quota, err
		}
	}
	return quota, nil
}

func update_quota(units *int, quota_reset_ts time.Time) {
	params := scout_db.Update_quota_params{Quota: int64(*units), QuotaReset: quota_reset_ts}
	err := queries.Update_quota(ctx, params)
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
