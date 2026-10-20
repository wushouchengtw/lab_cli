package helper

import buildbucketpb "go.chromium.org/luci/buildbucket/proto"

func ToBuildStatus(status buildbucketpb.Status) string {
	switch status {
	case buildbucketpb.Status_STATUS_UNSPECIFIED:
		return "UNSPECIFIED"
	case buildbucketpb.Status_SCHEDULED:
		return "SCHEDULED"
	case buildbucketpb.Status_STARTED:
		return "STARTED"
	case buildbucketpb.Status_ENDED_MASK:
		return "ENDED"
	case buildbucketpb.Status_SUCCESS:
		return "SUCCESS"
	case buildbucketpb.Status_FAILURE:
		return "FAILURE"
	case buildbucketpb.Status_INFRA_FAILURE:
		return "INFRA_FAILURE"
	case buildbucketpb.Status_CANCELED:
		return "CANCELED"
	default:
		return "UNKNOWN"
	}
}
