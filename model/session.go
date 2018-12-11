package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Session ... datatype for the ruleset
type Session struct {
	ID             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	IsActive       bool          `json:"isActive,omitempty"`
	StartDate      time.Time     `json:"startDate,omitempty"`
	DurationActive int64         `json:"durationActive,omitempty"`
	RulesetID      string        `json:"ruleset,omitempty"`
	Players        []*Player     `json:"players,omitempty"`
}

// NewSession ... creates a new session based ona ruleset
func NewSession(startDate time.Time, durationActive int64, ruleset string) *Session {
	return &Session{
		ID:             bson.NewObjectId(),
		StartDate:      startDate,
		DurationActive: durationActive,
		RulesetID:      ruleset,
		IsActive:       true,
	}
}

// AddPlayer .. adds player to an existing session
func (s *Session) AddPlayer(player *Player) {

	s.Players = append(s.Players, player)
}

// Player ... struct to keep track of player data
type Player struct {
	Name  string                 `json:"name,omitempty"`
	Stats map[string]interface{} `json:"stats,omitempty"`
}

// NewPlayer .. creates a new player
func NewPlayer(name string) *Player {
	return &Player{
		Name:  name,
		Stats: make(map[string]interface{}),
	}
}
