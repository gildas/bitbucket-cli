package common

import (
	"os"
	"strings"
)

func IsWSL() bool {
	buffer, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	version := strings.ToLower(string(buffer))
	return strings.Contains(version, "microsoft") || strings.Contains(version, "wsl")
}
