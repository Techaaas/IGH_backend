package main

func main1() {
	var db = database{nil}
	db.connector()
	db.dropTables()
	//db.createDifferenceTable()
	//db.createInfoTable()
	defer db.db.Close()
}
