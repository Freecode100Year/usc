package broker

import (
	"fmt"
	"strings"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

// NetworkBroker mediates all outbound egress requests at runtime.
type NetworkBroker struct {
	allowedCaps capability.Set
}

// NewNetworkBroker creates a network mediation broker.
func NewNetworkBroker(allowed capability.Set) *NetworkBroker {
	return &NetworkBroker{allowedCaps: allowed}
}

// AuthorizeRequest evaluates whether an HTTP request conforms to granted capabilities.
func (nb *NetworkBroker) AuthorizeRequest(method, host, path string) error {
	reqCap := capability.Capability{
		Kind:    capability.KindNetHTTP,
		Actions: []string{strings.ToUpper(method)},
		Resource: capability.Resource{
			Scheme: "https",
			Host:   capability.HostSpec{Type: "EXACT", Value: host},
			Path:   path,
		},
	}

	for _, granted := range nb.allowedCaps {
		if capability.Subsumes(granted, reqCap) {
			return nil
		}
	}
	return fmt.Errorf("USC-E7008: INV-8 violation: outbound network call blocked to %s%s [%s]", host, path, method)
}
