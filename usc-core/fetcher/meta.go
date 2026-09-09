package fetcher

// SourceType represents the recognized remote or local source format.
type SourceType string

const (
	SourceTypeLocalDir SourceType = "local_dir"
	SourceTypeClawHub  SourceType = "clawhub"
	SourceTypeGitHub   SourceType = "github"
	SourceTypeGeneric  SourceType = "generic_url"
)

// SkillMetadata encapsulates extracted declarative metadata for a skill.
type SkillMetadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Homepage    string   `json:"homepage,omitempty"`
	SourceURL   string   `json:"source_url"`
	Endpoints   []string `json:"endpoints"`
	EnvVars     []string `json:"env_vars"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	RawDoc      string   `json:"raw_doc,omitempty"`
	LocalPath   string   `json:"local_path"`
}
