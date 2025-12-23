package web

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseISO8601Datetime(t *testing.T) {
	tt := []struct {
		name    string
		date    string
		want    time.Time
		wantErr bool
	}{
		{
			name: "empty date",
			date: "",
			want: time.Time{},
		},
		{
			name:    "invalid date",
			date:    "invalid",
			wantErr: true,
		},
		{
			name: "valid date",
			date: "2021-01-01T00:00:00Z",
			want: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			got, err := parseISO8601Datetime(tc.date)
			if tc.wantErr {
				rq.Error(err)
			} else {
				rq.NoError(err)
				rq.Equal(tc.want, got)
			}
		})
	}
}
