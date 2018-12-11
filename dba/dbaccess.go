package dba

import (
	"fmt"

	"github.com/atla/athena/model"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DBAccess ... datatype for the ruleset
type DBAccess struct {
	Session *mgo.Session
}

// NewDBAccess ... creates a new ruleset
func NewDBAccess(session *mgo.Session) *DBAccess {
	return &DBAccess{
		Session: session,
	}
}

// GetAllRulesets .. asd
func (dba *DBAccess) GetAllRulesets() []model.Ruleset {

	c := dba.Session.DB("athenadb").C("ruleset")

	var result []model.Ruleset

	if err := c.Find(bson.M{}).All(&result); err == nil {
		return result
	} else {
	}

	return nil
}

// GetAllSessions .. asd
func (dba *DBAccess) GetAllSessions() []model.Session {

	c := dba.Session.DB("athenadb").C("session")

	var result []model.Session

	if err := c.Find(bson.M{}).All(&result); err == nil {
		return result
	} else {
	}

	return nil
}

// UpdateSession ... updates existing session
func (dba *DBAccess) UpdateSession(id string, update *model.Session) error {

	c := dba.Session.DB("athenadb").C("session")
	bsonID := bson.ObjectIdHex(id)

	// find existing session and update
	if session, err := dba.GetSessionByID(id); err != nil {

		// only update player stats
		for _, player := range update.Players {
			for _, p := range session.Players {
				if p.Name == player.Name {

					for key, value := range player.Stats {
						// override stat with new value
						p.Stats[key] = value
					}

				}
			}
		}

		if err := c.UpdateId(bsonID, session); err != nil {
			fmt.Println("Error Updating session with id: " + id)
			return err
		}
	} else {
		return err
	}

	return nil
}

// GetRulesetByID .. asd
func (dba *DBAccess) GetRulesetByID(id string) (*model.Ruleset, error) {

	c := dba.Session.DB("athenadb").C("ruleset")
	var result model.Ruleset

	bsonID := bson.ObjectIdHex(id)

	if err := c.Find(bson.M{"_id": bsonID}).One(&result); err == nil {
		return &result, nil
	} else {
		return nil, err
	}
}

// GetSessionByID .. returns a single session or error
func (dba *DBAccess) GetSessionByID(id string) (*model.Session, error) {

	c := dba.Session.DB("athenadb").C("session")
	var result model.Session

	bsonID := bson.ObjectIdHex(id)

	if err := c.Find(bson.M{"_id": bsonID}).One(&result); err == nil {
		return &result, nil
	} else {
		return nil, err
	}
}

// WriteSessionToDB ... asd
func (dba *DBAccess) WriteSessionToDB(session *model.Session) {
	sessionCollection := dba.Session.DB("athenadb").C("session")

	if err := sessionCollection.Insert(session); err != nil {
		fmt.Println("Error inserting session " + err.Error())
	}
}

// WriteRulesetToDB ... asd
func (dba *DBAccess) WriteRulesetToDB(ruleset *model.Ruleset) {
	rulesetCollection := dba.Session.DB("athenadb").C("ruleset")

	if err := rulesetCollection.Insert(ruleset); err != nil {
		fmt.Println("Error inserting instance")
	}
}

// UpdateRuleset ... asd
func (dba *DBAccess) UpdateRuleset(ruleset *model.Ruleset) error {
	rulesetCollection := dba.Session.DB("athenadb").C("ruleset")

	if err := rulesetCollection.UpdateId(ruleset.ID, ruleset); err != nil {
		fmt.Println("Error Updating ruleset with id: " + ruleset.ID)
		return err
	}
	return nil
}
