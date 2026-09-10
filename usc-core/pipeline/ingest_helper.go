package pipeline

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Freecode100Year/usc/usc-core/capability"
	"github.com/Freecode100Year/usc/usc-core/intent"
)

var (
	urlPattern = regexp.MustCompile(`https?://([a-zA-Z0-9\-\.]+)(:[0-9]+)?(/[^"'\s\)\]]*)?`)
	fsPattern  = regexp.MustCompile(`(?:open|read|cat|file|Path)\s*\(\s*['"](/etc/[^'"]+|~/[^'"]+|\.[^'"]+)['"]`)
)

type ingestedMeta struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Endpoints   []string `json:"endpoints"`
	EnvVars     []string `json:"env_vars"`
}

func extractIngestMetadata(skillName, sourcePath string) ([]intent.DeclaredStep, []intent.ObservedFact, string) {
	steps, desc := scanDeclaredMetadata(skillName, sourcePath)
	facts := scanObservedFacts(sourcePath)
	if len(facts) == 0 {
		facts = buildFactsFromSteps(steps)
	}
	if len(steps) == 0 {
		steps = buildStepsFromFacts(facts)
	}
	return steps, facts, desc
}

func scanDeclaredMetadata(skillName, sourcePath string) ([]intent.DeclaredStep, string) {
	if meta, ok := loadJSONMeta(sourcePath); ok {
		return buildStepsFromMeta(meta), meta.Description
	}
	if mdPath := filepath.Join(sourcePath, "SKILL.md"); fileExists(mdPath) {
		return scanSkillMD(mdPath, skillName)
	}
	return defaultSteps(skillName), fmt.Sprintf("Autonomous skill %s", skillName)
}

func loadJSONMeta(sourcePath string) (*ingestedMeta, bool) {
	metaFile := filepath.Join(sourcePath, "metadata.json")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, false
	}
	var meta ingestedMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, false
	}
	return &meta, true
}

func scanSkillMD(path, skillName string) ([]intent.DeclaredStep, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultSteps(skillName), skillName
	}
	content := string(data)
	urls := extractUniqueHosts(content)
	desc := extractMDDescription(content, skillName)
	if len(urls) == 0 {
		return defaultSteps(skillName), desc
	}
	var steps []intent.DeclaredStep
	for i, u := range urls {
		steps = append(steps, intent.DeclaredStep{
			StepID:      fmt.Sprintf("step_%d", i+1),
			Description: fmt.Sprintf("Query %s", u),
			Action:      "http.get",
			Target:      u,
		})
	}
	return steps, desc
}

func scanObservedFacts(sourcePath string) []intent.ObservedFact {
	var facts []intent.ObservedFact
	_ = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || shouldSkipScan(path) {
			return nil
		}
		rel, _ := filepath.Rel(sourcePath, path)
		facts = append(facts, scanSourceFile(path, rel)...)
		return nil
	})
	return facts
}

func scanSourceFile(filePath, relPath string) []intent.ObservedFact {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()
	var facts []intent.ObservedFact
	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		facts = append(facts, parseLineFacts(line, relPath, lineNo)...)
	}
	return facts
}

func parseLineFacts(line, relPath string, lineNo int) []intent.ObservedFact {
	var facts []intent.ObservedFact
	if m := urlPattern.FindStringSubmatch(line); len(m) > 1 {
		host := m[1]
		facts = append(facts, intent.ObservedFact{
			FactID:     fmt.Sprintf("obs_net_%s_%d", relPath, lineNo),
			SourceFile: relPath,
			LineNumber: lineNo,
			RawSnippet: strings.TrimSpace(line),
			Capability: capability.Capability{
				CapabilityID: fmt.Sprintf("cap_net_%s", host),
				Kind:         capability.KindNetHTTP,
				Actions:      []string{"GET"},
				Resource: capability.Resource{
					Scheme: "https",
					Host:   capability.HostSpec{Type: "EXACT", Value: host},
					Path:   "/*",
				},
			},
		})
	}
	facts = append(facts, parseFSFacts(line, relPath, lineNo)...)
	return facts
}

func parseFSFacts(line, relPath string, lineNo int) []intent.ObservedFact {
	var facts []intent.ObservedFact
	if m := fsPattern.FindStringSubmatch(line); len(m) > 1 {
		p := m[1]
		facts = append(facts, intent.ObservedFact{
			FactID:     fmt.Sprintf("obs_fs_%s_%d", relPath, lineNo),
			SourceFile: relPath,
			LineNumber: lineNo,
			RawSnippet: strings.TrimSpace(line),
			Capability: capability.Capability{
				CapabilityID: fmt.Sprintf("cap_fs_%s_%d", relPath, lineNo),
				Kind:         capability.KindFSRead,
				Actions:      []string{"READ"},
				Resource:     capability.Resource{Path: p},
			},
		})
	}
	return facts
}

func extractUniqueHosts(text string) []string {
	matches := urlPattern.FindAllStringSubmatch(text, -1)
	seen := make(map[string]bool)
	var hosts []string
	for _, m := range matches {
		if len(m) > 1 {
			h := m[1]
			if !seen[h] && !strings.Contains(h, "github.com") && !strings.Contains(h, "clawhub.ai") {
				seen[h] = true
				hosts = append(hosts, h)
			}
		}
	}
	return hosts
}

func extractMDDescription(text, fallback string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "description:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		}
		if strings.HasPrefix(line, "> ") {
			return strings.TrimPrefix(line, "> ")
		}
	}
	return fmt.Sprintf("Skill %s", fallback)
}

func buildStepsFromMeta(meta *ingestedMeta) []intent.DeclaredStep {
	var steps []intent.DeclaredStep
	for i, ep := range meta.Endpoints {
		steps = append(steps, intent.DeclaredStep{
			StepID:      fmt.Sprintf("step_%d", i+1),
			Description: fmt.Sprintf("Access %s", ep),
			Action:      "http.get",
			Target:      ep,
		})
	}
	return steps
}

func buildStepsFromFacts(facts []intent.ObservedFact) []intent.DeclaredStep {
	var steps []intent.DeclaredStep
	for i, f := range facts {
		steps = append(steps, intent.DeclaredStep{
			StepID:      fmt.Sprintf("step_%d", i+1),
			Description: "Step deduced from " + f.Capability.CapabilityID,
			Action:      f.Capability.Actions[0],
			Target:      f.Capability.Resource.Host.Value,
		})
	}
	return steps
}

func buildFactsFromSteps(steps []intent.DeclaredStep) []intent.ObservedFact {
	var facts []intent.ObservedFact
	for i, s := range steps {
		facts = append(facts, intent.ObservedFact{
			FactID:     fmt.Sprintf("obs_%d", i+1),
			SourceFile: "declared_spec",
			Capability: capability.Capability{
				CapabilityID: fmt.Sprintf("cap_step_%d", i+1),
				Kind:         capability.KindNetHTTP,
				Actions:      []string{"GET"},
				Resource: capability.Resource{
					Scheme: "https",
					Host:   capability.HostSpec{Type: "EXACT", Value: s.Target},
					Path:   "/*",
				},
			},
		})
	}
	return facts
}

func defaultSteps(skillName string) []intent.DeclaredStep {
	return []intent.DeclaredStep{
		{StepID: "main_step", Description: "Core skill execution", Action: "http.get", Target: "api.github.com"},
	}
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func shouldSkipScan(path string) bool {
	clean := filepath.ToSlash(path)
	return strings.Contains(clean, "/.git/") || strings.Contains(clean, "/node_modules/")
}
