package main

import (
	"fmt"

	dba "github.com/atla/athena/dba"

	"github.com/globalsign/mgo"
)

func main() {
	if session, err := mgo.Dial("athenadb"); err == nil {

		app := &App{
			DBAccess: dba.NewDBAccess(session),
		}

		app.CreateDatabase()
		app.SetupRoutes()
		app.Start()

	} else if err != nil {
		fmt.Println("Error connecting to mongo database: " + err.Error())
	}
}
