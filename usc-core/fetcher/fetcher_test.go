package fetcher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectSourceType(t *testing.T) {
	tests := []struct {
		in   string
		want SourceType
	}{
		{"https://clawhub.ai/thesentitrader/skills/us-stocks-analysis", SourceTypeClawHub},
		{"@thesentitrader/us-stocks-analysis", SourceTypeClawHub},
		{"https://github.com/Freecode100Year/usc", SourceTypeGitHub},
		{"./untrusted-skill", SourceTypeLocalDir},
	}
	for _, tt := range tests {
		if got := DetectSourceType(tt.in); got != tt.want {
			t.Errorf("DetectSourceType(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseClawHubSlug(t *testing.T) {
	owner, skill, ok := ParseClawHubSlug("https://clawhub.ai/thesentitrader/skills/us-stocks-analysis")
	if !ok || owner != "thesentitrader" || skill != "us-stocks-analysis" {
		t.Errorf("unexpected parse result: %s, %s, %v", owner, skill, ok)
	}
	owner2, skill2, ok2 := ParseClawHubSlug("@thesentitrader/us-stocks-analysis")
	if !ok2 || owner2 != "thesentitrader" || skill2 != "us-stocks-analysis" {
		t.Errorf("unexpected slug parse result: %s, %s, %v", owner2, skill2, ok2)
	}
}

func TestInspectLocalDir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "usc-fetch-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	meta, err := Fetch(tmp, tmp)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	if meta.Name != filepath.Base(tmp) {
		t.Errorf("expected name %s, got %s", filepath.Base(tmp), meta.Name)
	}
}
