package model

// SessionStart ... datatype for the ruleset
type SessionStart struct {
	RulesetID string   `json:"ruleset,omitempty"`
	Players   []string `json:"players,omitempty"`
}

// NewSessionStart ... creates a new session start object
func NewSessionStart(ruleset string, players []string) *SessionStart {
	return &SessionStart{
		RulesetID: ruleset,
		Players:   players,
	}
}
