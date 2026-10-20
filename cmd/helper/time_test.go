package helper

import (
	"testing"
)

func TestUnixToTaipeiTime(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		nanos    int32
		expected string
	}{
		{
			name:     "unix epoch",
			seconds:  0,
			nanos:    0,
			expected: "1970-01-01 08:00:00",
		},
		{
			name:     "specific timestamp",
			seconds:  1700000000,
			nanos:    500000000,
			expected: "2023-11-15 06:13:20",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := UnixToTaipeiTime(tc.seconds, tc.nanos)
			if got != tc.expected {
				t.Errorf("UnixToTaipeiTime(%d, %d) = %q, want %q", tc.seconds, tc.nanos, got, tc.expected)
			}
		})
	}
}
