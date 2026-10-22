package helper

import (
	"testing"

	buildbucketpb "go.chromium.org/luci/buildbucket/proto"
)

func TestToBuildStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   buildbucketpb.Status
		expected string
	}{
		{
			name:     "unspecified",
			status:   buildbucketpb.Status_STATUS_UNSPECIFIED,
			expected: "UNSPECIFIED",
		},
		{
			name:     "scheduled",
			status:   buildbucketpb.Status_SCHEDULED,
			expected: "SCHEDULED",
		},
		{
			name:     "started",
			status:   buildbucketpb.Status_STARTED,
			expected: "STARTED",
		},
		{
			name:     "ended_mask",
			status:   buildbucketpb.Status_ENDED_MASK,
			expected: "ENDED",
		},
		{
			name:     "success",
			status:   buildbucketpb.Status_SUCCESS,
			expected: "SUCCESS",
		},
		{
			name:     "failure",
			status:   buildbucketpb.Status_FAILURE,
			expected: "FAILURE",
		},
		{
			name:     "infra_failure",
			status:   buildbucketpb.Status_INFRA_FAILURE,
			expected: "INFRA_FAILURE",
		},
		{
			name:     "canceled",
			status:   buildbucketpb.Status_CANCELED,
			expected: "CANCELED",
		},
		{
			name:     "unknown",
			status:   buildbucketpb.Status(9999),
			expected: "UNKNOWN",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ToBuildStatus(tc.status)
			if got != tc.expected {
				t.Errorf("ToBuildStatus(%v) = %q, want %q", tc.status, got, tc.expected)
			}
		})
	}
}
