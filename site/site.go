package site

import (
	"os"
)

const ServiceAccountKeyPathEnv = "SERVICE_ACCOUNT_KEY_PATH"

func GetServiceAccountPath(injectKey ...string) string {
	if len(injectKey) > 0 && injectKey[0] != "" {
		return injectKey[0]
	}
	return os.Getenv(ServiceAccountKeyPathEnv)
}
