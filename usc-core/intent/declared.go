package intent

// DeclaredStep represents a step declared by the author.
type DeclaredStep struct {
	StepID      string   `json:"step_id"`
	Description string   `json:"description"`
	Action      string   `json:"action"`
	Target      string   `json:"target"`
	Parameters  []string `json:"parameters,omitempty"`
}

// DeclaredIntent encapsulates the author's claimed functionality (I_decl).
type DeclaredIntent struct {
	SkillName   string         `json:"skill_name"`
	Description string         `json:"description"`
	Steps       []DeclaredStep `json:"steps"`
	DeclaredEnv []string       `json:"declared_env,omitempty"`
}
