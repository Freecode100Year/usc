package broker

import (
	"os"
	"testing"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

func TestCredentialBroker(t *testing.T) {
	cb := NewCredentialBroker()
	cb.RegisterHandle("credential://github/token", "ghp_secret12345", "api.github.com")

	if _, err := cb.AuthorizeAndInject("credential://github/token", "evil.com"); err == nil {
		t.Fatalf("expected authorization failure on destination mismatch")
	}
	sec, err := cb.AuthorizeAndInject("credential://github/token", "api.github.com")
	if err != nil || sec != "ghp_secret12345" {
		t.Fatalf("expected successful injection, got %s, err: %v", sec, err)
	}
}

func TestNetworkAndFilesystemBrokers(t *testing.T) {
	allowed := capability.NewSet(capability.Capability{
		Kind:    capability.KindNetHTTP,
		Actions: []string{"GET"},
		Resource: capability.Resource{
			Host: capability.HostSpec{Type: "EXACT", Value: "api.weather.gov"},
			Path: "/*",
		},
	})
	nb := NewNetworkBroker(allowed)
	if err := nb.AuthorizeRequest("GET", "api.weather.gov", "/forecast"); err != nil {
		t.Fatalf("expected allowed network request, got: %v", err)
	}
	if err := nb.AuthorizeRequest("POST", "api.weather.gov", "/forecast"); err == nil {
		t.Fatalf("expected method mismatch block")
	}

	fb, _ := NewFilesystemBroker(os.TempDir(), true)
	if _, err := fb.AuthorizePath("../../../etc/passwd", false); err == nil {
		t.Fatalf("expected path traversal detection")
	}
}
