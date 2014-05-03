package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

type MyDb struct {
	gorp.DbMap
}

func (db *MyDb) addTestData() {
	// delete existing entries
	//  err := dbmap.TruncateTables()
	//  checkErr(err, "Truncate failed")

	// insert some seed data
	u1, _ := db.newUser("email", "a@b.c", "abcd1234", "Boo", "Daa")
	u2, _ := db.newUser("email", "b@c.d", "1234abcd", "Nee", "Naa")

	err := db.Insert(u1, u2)
	if err != nil {
		fmt.Println("failed to insert: ", err)
	}
}

func NewDb() *MyDb {
	db, err := sql.Open("postgres", "user=timcheng dbname=mila_dev sslmode=disable")
	checkErr(err, "sql.Open failed")
	dbmap := &MyDb{gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Connection{}, "connections").SetKeys(false, "user1_id", "user2_id")
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	dbmap.addTestData()
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
