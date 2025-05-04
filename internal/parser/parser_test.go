package parser

import (
	"reflect"
	"testing"
	"time"

	"YadroTest/internal/models"
)

func TestParseEventLine(t *testing.T) {
	testCases := []struct {
		name    string
		line    string
		want    models.Event
		wantErr bool
	}{
		{
			name: "Valid event with no extra params",
			line: "[09:30:01.005] 4 1",
			want: models.Event{
				Time:         parseTime(t, "09:30:01.005"),
				EventID:      4,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "Valid event with extra params",
			line: "[09:49:31.659] 5 1 1",
			want: models.Event{
				Time:         parseTime(t, "09:49:31.659"),
				EventID:      5,
				CompetitorID: 1,
				ExtraParams:  "1",
			},
			wantErr: false,
		},
		{
			name:    "Invalid time format",
			line:    "[9:30:01.5] 4 1",
			wantErr: true,
		},
		{
			name:    "Missing time brackets",
			line:    "09:30:01.005 4 1",
			wantErr: true,
		},
		{
			name:    "Missing event ID",
			line:    "[09:30:01.005]",
			wantErr: true,
		},
		{
			name:    "Invalid event ID",
			line:    "[09:30:01.005] x 1",
			wantErr: true,
		},
		{
			name:    "Invalid competitor ID",
			line:    "[09:30:01.005] 4 x",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseEventLine(tc.line)

			if (err != nil) != tc.wantErr {
				t.Errorf("ParseEventLine() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseEventLine() = %v, want %v", got, tc.want)
			}
		})
	}
}

func parseTime(t *testing.T, timeStr string) time.Time {
	result, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		t.Fatalf("Invalid test time: %s", timeStr)
	}
	return result
}
