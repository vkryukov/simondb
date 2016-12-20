// Database-related functions

package main

import (
	"log"

	"github.com/cznic/ql"
)

var _DB *ql.DB // Main database
const dbName = "./movies.db"

func mustCreateTable(v interface{}) {
	schema, err := ql.Schema(v, "", nil)
	if err != nil {
		log.Fatalf("Cannot create schema: %s\n", err)
	}
	if _, _, err := _DB.Execute(ql.NewRWCtx(), schema); err != nil {
		log.Fatalf("Cannot create schema: %s\n", err)
	}
}

func init() {
	var err error
	_DB, err = ql.OpenFile(dbName, &ql.Options{CanCreate: true})
	if err != nil {
		log.Fatalf("Cannot open database file %s: %s\n", dbName, err)
	}

	mustCreateTable(&Movie{})
	mustCreateTable(&Actor{})
	mustCreateTable(&MovieActor{})
}
