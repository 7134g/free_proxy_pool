package util

import "strings"

func FixScheme(s string) string {
	s = strings.ToLower(s)
	switch {
	case strings.Contains(s, "https"):
		return "https"
	case strings.Contains(s, "sock"):
		return "socks"
	default:
		return "http"
	}
}
