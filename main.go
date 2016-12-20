package main

const originalDBName = "./simondb.csv"

func main() {
	defer _DB.Close()
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
