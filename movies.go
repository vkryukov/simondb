// Structs and methods for managing movies and actors information

package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

// Movie holds parsed information about a movie
type Movie struct {
	ID int // Internal movie id

	// As recorded in database
	Title    string
	Studio   string
	Year     int
	Duration int
	Actors   []int // Internal actor IDs

	Kassette string // Stores kassete information in Simon's collection
	DVD      string // Stores DVD information in Simon's collection

	// Canonical information as per IMDB
	imdbID       string
	imdbTitle    string
	imdbStudio   string
	imdbYear     int
	imdbDuration int // Duration in minutes

	updated time.Time // last it was updated from IMDB
}

// Actor holds parsed information about actors in the movie
type Actor struct {
	ID     int // Internal actor id
	Name   string
	Movies []int // Internal movie IDs

	imdbID   string // IMDB id
	imdbName string // Canonical IMDB name

	updated time.Time // last it was updated from IMDB
}

// Movies hold all the movies in memory
type Movies []*Movie

// Actors hold all the actors in memory
type Actors []*Actor

// Serializing and de-serializing to/from disk

// ParseOriginalDB reads Simon's original file, and creates Actors and Movies collections
func ParseOriginalDB(filename string) (Movies, Actors, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var movies Movies
	actors := make(map[string]*Actor)
	maxActorID := 0
	for i, record := range records[1:] { // Skiping the headline
		// Fields are 0 - Title, 1 - Actors, 2 - Studio, 3 - Year, 4 - Duration (min), 5 - Kasette, 6 - DVD
		m := new(Movie)
		m.ID = i + 1
		m.Title = record[0]
		m.Studio = record[2]
		m.Year = parseYear(record[3])
		m.Duration = parseDuration(record[4])

		actorsNames := strings.Split(record[1], ",")
		for _, a := range actorsNames {
			a = strings.TrimSpace(a)
			if actor, ok := actors[a]; ok {
				// Actor is found, so we just need to add the MovieID to the list
				actor.Movies = append(actor.Movies, m.ID)
			} else {
				// Need to create a new actor
				maxActorID += 1
				actors[a] = new(Actor)
				actors[a].Name = a
				actors[a].ID = maxActorID
				actors[a].Movies = []int{m.ID}
			}
		}
	}

	var actorsList Actors
	for _, v := range actors {
		actorsList = append(actorsList, v)
	}
	return movies, actorsList, nil
}

// We expect properly parsed years, and return -1 on errors
func parseYear(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return int(i)
}

// We expect durationn in the form "XXX min.", and return -1 on errors
func parseDuration(s string) int {
	ss := strings.Split(s, " ")
	i, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		return -1
	}
	return int(i)
}
