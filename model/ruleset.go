package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Ruleset ... datatype for the ruleset
type Ruleset struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	DateCreated  time.Time     `json:"dateCreated,omitempty"`
	Game         string        `json:"game,omitempty"`
	Ruleset      string        `json:"ruleset,omitempty"`
	Stats        []*Stat       `json:"stats,omitempty"`
	Ranking      string        `json:"ranking,omitempty"`
	WinCondition string        `json:"winCondition,omitempty"`
	UsedInGames  int64         `json:"usedInGames,omitempty"`
}

// NewRuleset ... creates a new ruleset
func NewRuleset(game string, ruleset string) *Ruleset {
	return &Ruleset{
		ID:          bson.NewObjectId(),
		DateCreated: time.Now(),
		Game:        game,
		Ruleset:     ruleset,
	}
}

// AddStat ... adds new stat
func (r *Ruleset) AddStat(name string, label string, statType string, defaultValue interface{}) {

	stat := newStat(name, label, statType, defaultValue)
	r.Stats = append(r.Stats, stat)
}

func newStat(name string, label string, statType string, defaultValue interface{}) *Stat {
	return &Stat{
		Name:         name,
		Label:        label,
		StatType:     statType,
		DefaultValue: defaultValue,
	}
}

// Stat ... a single stat for a ruleset
type Stat struct {
	Name         string
	Label        string
	StatType     string
	DefaultValue interface{}
}
