package utils

import (
	"strings"
)

// SafeString returns a masked version of sensitive strings
func SafeString(value string, isSensitive bool) string {
	if !isSensitive || value == "" {
		return value
	}
	
	if len(value) <= 3 {
		return "***"
	}
	
	// Show first character and mask the rest
	return value[:1] + strings.Repeat("*", len(value)-1)
}

// IsIPAddress checks if a string looks like an IP address
func IsIPAddress(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return false
	}
	
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		// Simple check - not fully validating the number range
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}
	
	return true
}
