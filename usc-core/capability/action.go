package capability

import "strings"

const (
	ActionGet     = "GET"
	ActionPost    = "POST"
	ActionPut     = "PUT"
	ActionDelete  = "DELETE"
	ActionRead    = "READ"
	ActionWrite   = "WRITE"
	ActionExecute = "EXECUTE"
	ActionConnect = "CONNECT"
)

// ActionAllKnown contains the closed set of approved known actions.
var ActionAllKnown = []string{
	ActionGet, ActionPost, ActionPut, ActionDelete,
	ActionRead, ActionWrite, ActionExecute, ActionConnect,
}

// IsKnownAction verifies that an action belongs to the defended closed set.
func IsKnownAction(action string) bool {
	norm := strings.ToUpper(strings.TrimSpace(action))
	for _, known := range ActionAllKnown {
		if known == norm {
			return true
		}
	}
	return false
}

// ActionsSubsumes checks if parent action set subsumes child action set.
func ActionsSubsumes(parent, child []string) bool {
	parentMap := make(map[string]struct{}, len(parent))
	for _, a := range parent {
		parentMap[strings.ToUpper(strings.TrimSpace(a))] = struct{}{}
	}
	if _, hasStar := parentMap["*"]; hasStar {
		return true
	}
	for _, c := range child {
		cNorm := strings.ToUpper(strings.TrimSpace(c))
		if _, exists := parentMap[cNorm]; !exists {
			return false
		}
	}
	return true
}
