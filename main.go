package main

func main() {
	var db = database{nil}
	db.connector()
	db.dropTables()
	db.addDiffData([]string{"1", "2", "3"})
	//db.createDifferenceTable()
	//db.createInfoTable()
	defer db.db.Close()
}
