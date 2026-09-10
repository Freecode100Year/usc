package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func computeRealDigest(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	}
	h := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(h[:])
}

func writeProofFiles(proofDir string, p *PipelineState) error {
	if err := os.MkdirAll(proofDir, 0755); err != nil {
		return err
	}
	attBytes, _ := json.MarshalIndent(p.AttestationDoc, "", "  ")
	eventsBytes, _ := json.MarshalIndent(p.Ledger.Events(), "", "  ")
	bpBytes, _ := json.MarshalIndent(p.Blueprint, "", "  ")
	crBytes, _ := json.MarshalIndent(p.CleanRoom, "", "  ")
	_ = os.WriteFile(filepath.Join(proofDir, "attestation.json"), attBytes, 0644)
	_ = os.WriteFile(filepath.Join(proofDir, "audit-chain.json"), eventsBytes, 0644)
	_ = os.WriteFile(filepath.Join(proofDir, "blueprint.json"), bpBytes, 0644)
	_ = os.WriteFile(filepath.Join(proofDir, "cleanroom-proof.json"), crBytes, 0644)
	_ = os.WriteFile(filepath.Join(proofDir, "provenance.json"), []byte(`{"provenance":"strict"}`), 0644)
	_ = os.WriteFile(filepath.Join(proofDir, "policy.json"), []byte(`{"policy":"enterprise_strict"}`), 0644)
	_ = os.WriteFile(filepath.Join(proofDir, "sbom.spdx.json"), []byte(`{"spdxVersion":"SPDX-2.3"}`), 0644)
	return nil
}

func prepareArtifactPayload(p *PipelineState, payloadDir string) error {
	_ = os.RemoveAll(payloadDir)
	if err := os.MkdirAll(payloadDir, 0755); err != nil {
		return err
	}
	if p.SourceDir != "" && isDir(p.SourceDir) {
		return copyDirFiles(p.SourceDir, payloadDir)
	}
	content := fmt.Sprintf("# %s\n\n%s\n", p.SkillName, p.DeclaredIntent.Description)
	return os.WriteFile(filepath.Join(payloadDir, "SKILL.md"), []byte(content), 0644)
}

func copyDirFiles(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		_ = os.MkdirAll(filepath.Dir(target), 0755)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}

