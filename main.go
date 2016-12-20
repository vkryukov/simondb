package main

const originalDBName = "./simondb.csv"

func main() {
	movies, actors, _, err := ParseOriginalDB(originalDBName)
	if err != nil {
		panic(err)
	}
	if err = movies.WriteCSV("movies.csv"); err != nil {
		panic(err)
	}
	if err = actors.WriteCSV("actors.csv"); err != nil {
		panic(err)
	}
}
