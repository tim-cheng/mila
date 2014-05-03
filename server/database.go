package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

type MyDb struct {
	gorp.DbMap
}

func NewDb() *MyDb {
	db, err := sql.Open("postgres", "user=timcheng dbname=mila_dev sslmode=disable")
	checkErr(err, "sql.Open failed")
	dbmap := &MyDb{gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Connection{}, "connections").SetKeys(false, "user1_id", "user2_id")
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	dbmap.AddTableWithName(Comment{}, "comments").SetKeys(true, "Id")
	dbmap.AddTableWithName(Star{}, "stars").SetKeys(false, "post_id", "user_id")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
