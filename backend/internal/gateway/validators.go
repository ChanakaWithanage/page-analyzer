package gateway

import "regexp"

var urlRegex = regexp.MustCompile(`^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(:[0-9]+)?(\/\S*)?$`)

func isValidURL(raw string) bool {
	return urlRegex.MatchString(raw)
}
