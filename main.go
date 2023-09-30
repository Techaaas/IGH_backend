package main

func main() {
	var db = database{nil, false, false}
	db.connector()
	//db.dropTables()
	//db.createDifferenceTable()
	//db.createInfoTable()
	defer db.db.Close()
}
