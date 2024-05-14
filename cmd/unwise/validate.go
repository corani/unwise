package main

import "time"

func parseISO8601Datetime(date string) (time.Time, error) {
	if date == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
