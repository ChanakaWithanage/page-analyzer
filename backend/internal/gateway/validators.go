package gateway

import (
    "regexp"
    "strings"
)

var urlRegex = regexp.MustCompile(`^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(:[0-9]+)?(\/\S*)?$`)

func isValidURL(raw string) bool {
    if strings.Contains(raw, "127.0.0.1") || strings.Contains(raw, "localhost") {
        return true  // Allow local host for testing
    }
	return urlRegex.MatchString(raw)
}
