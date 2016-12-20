package main

const originalDBName = "./simondb.csv"

func main() {
	defer _DB.Close()
	initialImportExport()
}
