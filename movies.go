// Structs and methods for managing movies and actors information

package main

import (
	"encoding/csv"
	"os"
	"sort"
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
		m.Kassette = record[5]
		m.DVD = record[6]

		actorsNames := strings.Split(record[1], ",")
		for _, a := range actorsNames {
			a = strings.TrimSpace(a)
			if actor, ok := actors[a]; ok {
				// Actor is found, so we just need to add the MovieID to the list
				actor.Movies = append(actor.Movies, m.ID)
				m.Actors = append(m.Actors, actor.ID)
			} else {
				// Need to create a new actor
				maxActorID += 1
				actors[a] = new(Actor)
				actors[a].Name = a
				actors[a].ID = maxActorID
				actors[a].Movies = []int{m.ID}
				m.Actors = append(m.Actors, actors[a].ID)
			}
		}
		movies = append(movies, m)

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

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 03:04:05")
}

func (movies Movies) Len() int {
	return len(movies)
}

func (movies Movies) Swap(i, j int) {
	movies[i], movies[j] = movies[j], movies[i]
}

func (movies Movies) Less(i, j int) bool {
	return movies[i].Title < movies[j].Title
}

// WriteCSV writes movies to disk as a CSV file
func (movies Movies) WriteCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	// Header
	writer.Write([]string{
		"ID", "Title", "Studio", "Year", "Duration", "Actors",
		"Kassette", "DVD", "imdbID", "imdbTitle", "imdbStudio",
		"imdbYear", "imdbDuration", "updated",
	})
	if writer.Error() != nil {
		return writer.Error()
	}
	sort.Sort(movies)
	for _, m := range movies {
		var actors []string
		for _, a := range m.Actors {
			actors = append(actors, strconv.Itoa(a))
		}
		writer.Write([]string{
			strconv.Itoa(m.ID),
			m.Title,
			m.Studio,
			strconv.Itoa(m.Year),
			strconv.Itoa(m.Duration),
			strings.Join(actors, " : "),
			m.Kassette,
			m.DVD,
			m.imdbID,
			m.imdbTitle,
			m.imdbStudio,
			strconv.Itoa(m.imdbYear),
			strconv.Itoa(m.imdbDuration),
			formatTime(m.updated),
		})
		if writer.Error() != nil {
			return writer.Error()
		}
	}
	writer.Flush()
	return writer.Error()
}

func (actors Actors) Len() int {
	return len(actors)
}

func (actors Actors) Swap(i, j int) {
	actors[i], actors[j] = actors[j], actors[i]
}

func (actors Actors) Less(i, j int) bool {
	return actors[i].Name < actors[j].Name
}

// WriteCSV writes actors to disk as a CSV file
func (actors Actors) WriteCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	// Header
	writer.Write([]string{
		"ID", "Name", "Movies",
		"imdbID", "imdbName", "updated",
	})
	if writer.Error() != nil {
		return writer.Error()
	}
	sort.Sort(actors)
	for _, a := range actors {
		var movies []string
		for _, m := range a.Movies {
			movies = append(movies, strconv.Itoa(m))
		}
		writer.Write([]string{
			strconv.Itoa(a.ID),
			a.Name,
			strings.Join(movies, ","),
			a.imdbID,
			a.imdbName,
			formatTime(a.updated),
		})
		if writer.Error() != nil {
			return writer.Error()
		}
	}
	writer.Flush()
	return writer.Error()
}
