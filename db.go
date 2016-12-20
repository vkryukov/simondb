// Database-related functions

package main

import (
	"log"

	"github.com/cznic/ql"
)

var _DB *ql.DB // Main database
const dbName = "./movies.db"

var insertMovie, insertActor, insertMovieActor ql.List

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

	insertMovie = ql.MustCompile("insert into Movie values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);")
	insertActor = ql.MustCompile("insert into Actor values ($1, $2, $3, $4, $5);")
	insertMovieActor = ql.MustCompile("insert into MovieActor values ($1, $2);")

}

// SaveToDB save movies information to the database
func (movies Movies) SaveToDB() error {
	ctx := ql.NewRWCtx()
	if _, _, err := _DB.Run(ctx, "BEGIN TRANSACTION;"); err != nil {
		return err
	}
	for _, m := range movies {
		if _, _, err := _DB.Execute(
			ctx, insertMovie,
			m.ID, m.Title, m.Studio, m.Year, m.Duration, m.Kassette, m.DVD,
			m.IMDBID, m.IMDBTitle, m.IMDBStudio, m.IMDBYear, m.IMDBDuration, m.Updated,
		); err != nil {
			return err
		}
	}
	if _, _, err := _DB.Run(ctx, "COMMIT;"); err != nil {
		return err
	}
	if err := _DB.Flush(); err != nil {
		return err
	}
	return nil
}

// SaveToDB save actors information to the database
func (actors Actors) SaveToDB() error {
	ctx := ql.NewRWCtx()
	if _, _, err := _DB.Run(ctx, "BEGIN TRANSACTION;"); err != nil {
		return err
	}
	for _, a := range actors {
		if _, _, err := _DB.Execute(
			ctx, insertActor,
			a.ID, a.Name, a.IMDBID, a.IMDBName, a.Updated,
		); err != nil {
			return err
		}
	}
	if _, _, err := _DB.Run(ctx, "COMMIT;"); err != nil {
		return err
	}
	if err := _DB.Flush(); err != nil {
		return err
	}
	return nil
}

// SaveToDB save actors information to the database
func (moviesActors MoviesActors) SaveToDB() error {
	ctx := ql.NewRWCtx()
	if _, _, err := _DB.Run(ctx, "BEGIN TRANSACTION;"); err != nil {
		return err
	}
	for _, ma := range moviesActors {
		if _, _, err := _DB.Execute(
			ctx, insertMovieActor,
			ma.MovieID, ma.ActorID,
		); err != nil {
			return err
		}
	}
	if _, _, err := _DB.Run(ctx, "COMMIT;"); err != nil {
		return err
	}
	if err := _DB.Flush(); err != nil {
		return err
	}
	return nil
}

func initialImportExport() {
	movies, actors, moviesActors, err := ParseOriginalDB(originalDBName)
	if err != nil {
		log.Fatalf("Error parsing original file: %s\n", err)
	}
	if err = movies.WriteCSV("movies.csv"); err != nil {
		log.Fatalf("Cannot output movies.csv: %s\n", err)
	}
	if err = actors.WriteCSV("actors.csv"); err != nil {
		log.Fatalf("Cannot output actors.csv: %s\n", err)
	}
	if err = movies.SaveToDB(); err != nil {
		log.Fatalf("Cannot write movies to DB: %s\n", err)
	}
	if err = actors.SaveToDB(); err != nil {
		log.Fatalf("Cannot write actors to DB: %s\n", err)
	}
	if err = moviesActors.SaveToDB(); err != nil {
		log.Fatalf("Cannot write moviesActors to DB: %s\n", err)
	}
}
