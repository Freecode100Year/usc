package capability

// Kind represents the formal classification of a capability.
type Kind string

const (
	KindNetHTTP      Kind = "net.http"
	KindNetEgress    Kind = "net.egress"
	KindFSRead       Kind = "fs.read"
	KindFSWrite      Kind = "fs.write"
	KindSysExec      Kind = "sys.exec"
	KindSecretAccess Kind = "secret.access"
	KindEnvRead      Kind = "env.read"
)

// HostSpec defines the target host specification.
type HostSpec struct {
	Type  string `json:"type"`  // "EXACT", "WILDCARD"
	Value string `json:"value"` // e.g. "api.github.com", "*.github.com"
}

// Resource defines the precise target of a capability.
type Resource struct {
	Scheme string   `json:"scheme,omitempty"`
	Host   HostSpec `json:"host"`
	Port   int      `json:"port,omitempty"`
	Path   string   `json:"path,omitempty"`
}

// Constraints defines hard numerical and policy limits for capability execution.
type Constraints struct {
	AllowRedirect      bool  `json:"allow_redirect"`
	MaxResponseBytes   int64 `json:"max_response_bytes"`
	RateLimitPerMinute int   `json:"rate_limit_per_minute"`
}

// Justification binds a capability to a declared intent step.
type Justification struct {
	IntentStep string `json:"intent_step"`
	RequiredBy string `json:"required_by"`
}

// Capability represents the formal zero-trust capability contract structure.
type Capability struct {
	CapabilityID  string        `json:"capability_id"`
	Kind          Kind          `json:"kind"`
	Actions       []string      `json:"actions"`
	Resource      Resource      `json:"resource"`
	Constraints   Constraints   `json:"constraints"`
	Justification Justification `json:"justification"`
}
