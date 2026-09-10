package attestation

import (
	"archive/zip"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ComputeDirSHA256 computes a deterministic cryptographic hash of a directory tree.
func ComputeDirSHA256(dirPath string) (string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(p string, fi os.FileInfo, err error) error {
		if err == nil && !fi.IsDir() {
			files = append(files, p)
		}
		return err
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	h := sha256.New()
	for _, f := range files {
		rel, _ := filepath.Rel(dirPath, f)
		h.Write([]byte(filepath.ToSlash(rel) + "\n"))
		data, err := os.ReadFile(f)
		if err != nil {
			return "", err
		}
		h.Write(data)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SignAndPackageArtifact hashes the payload, signs the attestation, and packages the .usc zip container.
func SignAndPackageArtifact(destZip string, att *Attestation, priv ed25519.PrivateKey, keyID, proofDir, payloadDir string) error {
	payloadHash, err := ComputeDirSHA256(payloadDir)
	if err != nil {
		return fmt.Errorf("failed to hash payload: %w", err)
	}
	att.Artifact.SHA256 = "sha256:" + payloadHash
	if err := SignAttestation(att, priv, keyID); err != nil {
		return err
	}
	return PackageArtifact(destZip, att, proofDir, payloadDir)
}

// PackageArtifact packages attestation, proof, and payload into a .usc zip container.
func PackageArtifact(destZip string, att *Attestation, proofDir, payloadDir string) error {
	if err := os.MkdirAll(filepath.Dir(destZip), 0755); err != nil {
		return err
	}
	zf, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zf.Close()
	w := zip.NewWriter(zf)
	defer w.Close()
	if err := addAttestationToZip(w, att); err != nil {
		return err
	}
	if err := addDirToZip(w, proofDir, "proof"); err != nil {
		return err
	}
	return addDirToZip(w, payloadDir, "payload")
}

func addAttestationToZip(w *zip.Writer, att *Attestation) error {
	f, err := w.Create("attestation.json")
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(att, "", "  ")
	_, err = f.Write(data)
	return err
}

func addDirToZip(w *zip.Writer, srcDir, prefix string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(srcDir, path)
		zipPath := filepath.ToSlash(filepath.Join(prefix, rel))
		zf, err := w.Create(zipPath)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = zf.Write(data)
		return err
	})
}

// UnpackArtifact uncompresses a .usc file safely, protecting against zip slip attacks.
func UnpackArtifact(zipPath, destDir string) (*Attestation, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open .usc artifact: %w", err)
	}
	defer r.Close()
	for _, f := range r.File {
		if err := extractZipFile(f, destDir); err != nil {
			return nil, err
		}
	}
	return loadAttestation(filepath.Join(destDir, "attestation.json"))
}

func extractZipFile(f *zip.File, destDir string) error {
	cleanName := filepath.Clean(f.Name)
	if strings.HasPrefix(cleanName, "..") || filepath.IsAbs(cleanName) {
		return fmt.Errorf("USC-E7048: illegal zip path traversal %s", f.Name)
	}
	targetPath := filepath.Join(destDir, cleanName)
	if f.FileInfo().IsDir() {
		return os.MkdirAll(targetPath, 0755)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, rc)
	return err
}

func loadAttestation(path string) (*Attestation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("missing attestation.json: %w", err)
	}
	var att Attestation
	if err := json.Unmarshal(data, &att); err != nil {
		return nil, fmt.Errorf("corrupted attestation.json: %w", err)
	}
	return &att, nil
}
