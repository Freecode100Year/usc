package capability

import (
	"path"
	"strings"
)

// SegmentType distinguishes literal segments from wildcard matchers.
type SegmentType int

const (
	SegmentLiteral SegmentType = iota
	SegmentWildcard
)

// PathSegment represents a parsed AST node of a URL or filesystem path.
type PathSegment struct {
	Type  SegmentType
	Value string
}

// ParsePathSegments parses and canonicalizes a path into Segment AST.
func ParsePathSegments(p string) []PathSegment {
	clean := path.Clean(strings.TrimSpace(p))
	if clean == "." || clean == "/" || clean == "" {
		return []PathSegment{}
	}
	parts := strings.Split(strings.Trim(clean, "/"), "/")
	segments := make([]PathSegment, 0, len(parts))
	for _, part := range parts {
		if part == "*" {
			segments = append(segments, PathSegment{Type: SegmentWildcard, Value: "*"})
		} else {
			segments = append(segments, PathSegment{Type: SegmentLiteral, Value: part})
		}
	}
	return segments
}

// PathSubsumes checks if parent path pattern subsumes child path pattern.
func PathSubsumes(parent, child string) bool {
	pSegs := ParsePathSegments(parent)
	cSegs := ParsePathSegments(child)
	if len(pSegs) == 0 {
		return parent == "*" || len(cSegs) == 0
	}
	if len(cSegs) < len(pSegs) {
		return false
	}
	for i := 0; i < len(pSegs); i++ {
		isLast := i == len(pSegs)-1
		if pSegs[i].Type == SegmentWildcard && isLast {
			return true
		}
		if pSegs[i].Type == SegmentWildcard {
			continue
		}
		if pSegs[i].Value != cSegs[i].Value {
			return false
		}
	}
	return len(pSegs) == len(cSegs)
}

// HostSubsumes determines if parent host pattern subsumes child host.
func HostSubsumes(parent, child HostSpec) bool {
	if parent.Type == "WILDCARD" && parent.Value == "*" {
		return true
	}
	if parent.Type == "WILDCARD" && strings.HasPrefix(parent.Value, "*.") {
		suffix := strings.TrimPrefix(parent.Value, "*.")
		return strings.HasSuffix(child.Value, suffix)
	}
	return parent.Value == child.Value
}

// ResourceSubsumes checks whether parent resource encompasses child resource.
func ResourceSubsumes(parent, child Resource) bool {
	if parent.Scheme != "" && parent.Scheme != "*" && parent.Scheme != child.Scheme {
		return false
	}
	if !HostSubsumes(parent.Host, child.Host) {
		return false
	}
	if parent.Port != 0 && child.Port != 0 && parent.Port != child.Port {
		return false
	}
	return PathSubsumes(parent.Path, child.Path)
}
