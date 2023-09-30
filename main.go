package main

var db database

func main1() {
	db.connector()
	//db.dropTables()
	//db.createDifferenceTable()
	//db.createInfoTable()
	defer db.db.Close()
}
