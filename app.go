package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dba "github.com/atla/athena/dba"
	model "github.com/atla/athena/model"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

// App ... main application structure
type App struct {
	DBAccess *dba.DBAccess
	Router   *mux.Router
}

// GetSessions ... returns the list of sessions
func (app *App) GetSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	result := app.DBAccess.GetAllSessions()
	now := time.Now()

	// update all IsActive states based on duration active
	for _, session := range result {
		session.IsActive = now.Sub(session.StartDate).Minutes() > float64(session.DurationActive)
	}

	json.NewEncoder(w).Encode(result)
}

// GetRulesets ... returns the list of rulesets
func (app *App) GetRulesets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var result = app.DBAccess.GetAllRulesets()
	json.NewEncoder(w).Encode(result)
}

// StartSession ... starts a new session based on the send parameters
func (app *App) CreateRuleset(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var ruleset model.Ruleset
	_ = json.NewDecoder(r.Body).Decode(&ruleset)

	// update newly generated ruleset
	ruleset.ID = bson.NewObjectId()
	ruleset.DateCreated = time.Now()

	app.DBAccess.WriteRulesetToDB(&ruleset)

	respondWithJSON(w, http.StatusOK, &ruleset)
}

// responds to the request with the given code and payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// StartSession ... starts a new session based on the send parameters
func (app *App) StartSession(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var sessionStart model.SessionStart
	_ = json.NewDecoder(r.Body).Decode(&sessionStart)

	session, err := app.StartSessionWithParams(&sessionStart)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "no such ruleset id")
	} else {
		respondWithJSON(w, http.StatusOK, session)
	}
}

// StartSessionWithParams ... creates a new Session based on the parameters of SessionStart
func (app *App) StartSessionWithParams(sessionStart *model.SessionStart) (*model.Session, error) {

	session := model.NewSession(time.Now(), 5, sessionStart.RulesetID)
	ruleset, err := app.DBAccess.GetRulesetByID(sessionStart.RulesetID)

	if err != nil {
		return nil, err
	}

	for _, playerName := range sessionStart.Players {

		player := model.NewPlayer(playerName)

		for _, stat := range ruleset.Stats {
			player.Stats[stat.Name] = stat.DefaultValue
		}

		session.AddPlayer(player)
	}

	// store session in db
	app.DBAccess.WriteSessionToDB(session)

	// track how many time this ruleset was played and store in db
	ruleset.UsedInGames++
	app.DBAccess.UpdateRuleset(ruleset)

	return session, nil
}

// GetRuleset ... returns the list of rulesets
func (app *App) GetRuleset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	var id = params["id"]

	ruleset, err := app.DBAccess.GetRulesetByID(id)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "no such ruleset id")
	} else {
		respondWithJSON(w, http.StatusOK, ruleset)
	}

}

// UpdateSession ... returns the list of rulesets
func (app *App) UpdateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	var id = params["id"]
	var sessionUpdate model.Session
	_ = json.NewDecoder(r.Body).Decode(&sessionUpdate)

	err := app.DBAccess.UpdateSession(id, &sessionUpdate)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "no such session id")
	} else {
		respondWithJSON(w, http.StatusOK, "session")
	}

}

// CreateDatabase ... creates initial database value
func (app *App) CreateDatabase() {

	// reset database
	app.DBAccess.Session.DB("athenadb").C("ruleset").DropCollection()
	app.DBAccess.Session.DB("athenadb").C("ruleset").Create(&mgo.CollectionInfo{})
	app.DBAccess.Session.DB("athenadb").C("session").DropCollection()
	app.DBAccess.Session.DB("athenadb").C("session").Create(&mgo.CollectionInfo{})

	munchkinRules := model.NewRuleset("Munchkin", "Classic")
	munchkinRules.Ranking = "{playerLevel} descending"
	munchkinRules.WinCondition = "{playerLevel} > 9"

	munchkinRules.AddStat("playerLevel", "Player Level", "int", 1)
	munchkinRules.AddStat("itemLevel", "Item Level", "int", 0)
	munchkinRules.AddStat("battleStrength", "Battle Level", "int", 0)

	duration, _ := time.ParseDuration("2h")

	munchkinSession := model.NewSession(time.Now(), duration.Nanoseconds(), munchkinRules.ID.Hex())
	munchkinSession.AddPlayer(model.NewPlayer("atla"))
	munchkinSession.AddPlayer(model.NewPlayer("claudia"))
	munchkinSession.AddPlayer(model.NewPlayer("daniel"))

	// store sample data
	app.DBAccess.WriteRulesetToDB(munchkinRules)
	app.DBAccess.WriteSessionToDB(munchkinSession)
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

// SetupRoutes ... Configures the routes
func (app *App) SetupRoutes() {

	var routes = Routes{
		Route{
			"Get all Rulesets",
			"GET",
			"/api/ruleset",
			app.GetRulesets,
		},
		Route{
			"Craete Ruleset",
			"POST",
			"/api/ruleset",
			app.CreateRuleset,
		},
		Route{
			"Get Ruleset",
			"GET",
			"/api/ruleset/{id}",
			app.GetRuleset,
		},
		Route{
			"Get all Session",
			"GET",
			"/api/session",
			app.GetSessions,
		},
		Route{
			"Start Session",
			"POST",
			"/api/startsession",
			app.StartSession,
		},
		Route{
			"Update Session",
			"PATCH",
			"/api/session/{id}",
			app.UpdateSession,
		},
	}

	// wrap all routes in logger

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	app.Router = router
	// also setup static serving
	app.Router.PathPrefix("/app").Handler(http.FileServer(http.Dir("dist/")))
}

// Start ... starts the server
func (app *App) Start() {
	fmt.Println("Server is running")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", app.Router))
}
