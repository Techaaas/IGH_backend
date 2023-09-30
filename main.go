package main

func main() {
	var db = database{nil}
	db.connector()
	db.dropTables()
	//db.createDifferenceTable()
	//db.createInfoTable()
	defer db.db.Close()
}
