package fetcher

import (
	"net/url"
	"strings"
)

// IsRemoteSource checks whether the source path is a remote URL or remote slug.
func IsRemoteSource(src string) bool {
	s := strings.TrimSpace(src)
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return true
	}
	if strings.HasPrefix(s, "github.com/") || strings.HasPrefix(s, "clawhub.ai/") {
		return true
	}
	if strings.HasPrefix(s, "@") && strings.Contains(s, "/") {
		return true
	}
	return false
}

// DetectSourceType classifies the source target into a known type.
func DetectSourceType(src string) SourceType {
	s := strings.ToLower(strings.TrimSpace(src))
	if strings.Contains(s, "clawhub.ai") || (strings.HasPrefix(s, "@") && strings.Contains(s, "/")) {
		return SourceTypeClawHub
	}
	if strings.Contains(s, "github.com") {
		return SourceTypeGitHub
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return SourceTypeGeneric
	}
	return SourceTypeLocalDir
}

// ParseClawHubSlug extracts the owner and skill slug from a ClawHub URL or slug.
func ParseClawHubSlug(src string) (string, string, bool) {
	s := strings.TrimSpace(src)
	if strings.HasPrefix(s, "@") {
		parts := strings.Split(strings.TrimPrefix(s, "@"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], true
		}
	}
	u, err := url.Parse(s)
	if err == nil && strings.Contains(u.Host, "clawhub.ai") {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) >= 3 && parts[1] == "skills" {
			return parts[0], parts[2], true
		}
		if len(parts) == 2 {
			return parts[0], parts[1], true
		}
	}
	return "", "", false
}
