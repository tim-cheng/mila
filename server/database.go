package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

func addTestData(dbmap *gorp.DbMap) {
	// delete existing entries
	//  err := dbmap.TruncateTables()
	//  checkErr(err, "Truncate failed")

	// insert some seed data
	u1 := newUser("email", "a@b.c", "Boo", "Daa")
	u2 := newUser("email", "b@c.d", "Nee", "Naa")
	err := dbmap.Insert(&u1, &u2)
	checkErr(err, "insert failed")
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "user=timcheng dbname=mila_dev sslmode=disable")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	addTestData(dbmap)

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
