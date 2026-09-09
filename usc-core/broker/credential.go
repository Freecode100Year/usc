package broker

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrSecretNotFound   = errors.New("USC-E7010: credential handle not found in vault")
	ErrIllegalSecretEnv = errors.New("USC-E7011: INV-5 violation: secret material in process environment")
)

// CredentialBroker enforces INV-5: Secrets as Capabilities.
type CredentialBroker struct {
	mu     sync.RWMutex
	vault  map[string]string // handle -> raw secret
	scopes map[string]string // handle -> allowed host destination
}

// NewCredentialBroker initializes a secure broker instance.
func NewCredentialBroker() *CredentialBroker {
	return &CredentialBroker{
		vault:  make(map[string]string),
		scopes: make(map[string]string),
	}
}

// RegisterHandle maps a capability handle to a secret and destination scope.
func (b *CredentialBroker) RegisterHandle(handle, rawSecret, allowedHost string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.vault[handle] = rawSecret
	b.scopes[handle] = allowedHost
}

// AuthorizeAndInject verifies destination host before injecting secret into request header.
func (b *CredentialBroker) AuthorizeAndInject(handle, destHost string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	expectedHost, ok := b.scopes[handle]
	if !ok {
		return "", ErrSecretNotFound
	}
	if !strings.EqualFold(expectedHost, destHost) {
		return "", errors.New("USC-E7012: credential destination mismatch: unauthorized declassification attempt")
	}
	return b.vault[handle], nil
}

// AuditEnvironment asserts that raw secret values are not present in ambient env.
func (b *CredentialBroker) AuditEnvironment(envVars []string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			val := parts[1]
			for _, secretVal := range b.vault {
				if len(secretVal) > 4 && strings.Contains(val, secretVal) {
					return ErrIllegalSecretEnv
				}
			}
		}
	}
	return nil
}
