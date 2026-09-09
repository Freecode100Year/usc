package broker

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// ProcessBroker mediates subprocess execution attempts according to granted policy.
type ProcessBroker struct {
	allowedBinaries map[string]struct{}
	allowExec       bool
}

// NewProcessBroker creates a process boundary mediator.
func NewProcessBroker(allowExec bool, allowedBinaries []string) *ProcessBroker {
	bins := make(map[string]struct{}, len(allowedBinaries))
	for _, b := range allowedBinaries {
		bins[strings.ToLower(filepath.Base(b))] = struct{}{}
	}
	return &ProcessBroker{
		allowedBinaries: bins,
		allowExec:       allowExec,
	}
}

// AuthorizeExecution verifies if an execve call is authorized.
func (pb *ProcessBroker) AuthorizeExecution(binaryPath string) error {
	if !pb.allowExec {
		return errors.New("USC-E7030: INV-8 violation: subprocess execution prohibited by policy")
	}
	base := strings.ToLower(filepath.Base(binaryPath))
	if _, ok := pb.allowedBinaries[base]; !ok {
		return fmt.Errorf("USC-E7031: subprocess execution denied for binary: %s", base)
	}
	return nil
}
