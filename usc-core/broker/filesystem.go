package broker

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrPathTraversal     = errors.New("USC-E7020: forbidden directory traversal detected")
	ErrOutOfSandboxScope = errors.New("USC-E7021: file operation outside authorized sandbox root")
)

// FilesystemBroker guards filesystem access against path traversal and unauthorized scopes.
type FilesystemBroker struct {
	sandboxRoot string
	readOnly    bool
}

// NewFilesystemBroker initializes a guarded filesystem broker.
func NewFilesystemBroker(sandboxRoot string, readOnly bool) (*FilesystemBroker, error) {
	abs, err := filepath.Abs(sandboxRoot)
	if err != nil {
		return nil, err
	}
	return &FilesystemBroker{
		sandboxRoot: filepath.Clean(abs),
		readOnly:    readOnly,
	}, nil
}

// AuthorizePath checks canonical path against directory traversal and sandbox boundary.
func (fb *FilesystemBroker) AuthorizePath(targetPath string, isWrite bool) (string, error) {
	if isWrite && fb.readOnly {
		return "", errors.New("USC-E7022: write rejected on read-only sandbox mount")
	}
	cleanRel := filepath.Clean(targetPath)
	if strings.HasPrefix(cleanRel, "..") || strings.Contains(cleanRel, "/../") || strings.Contains(cleanRel, "\\..\\") {
		return "", ErrPathTraversal
	}
	resolved := filepath.Clean(filepath.Join(fb.sandboxRoot, cleanRel))
	if !strings.HasPrefix(resolved, fb.sandboxRoot) {
		return "", ErrOutOfSandboxScope
	}
	return resolved, nil
}
