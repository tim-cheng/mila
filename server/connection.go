package main

import (
	"time"
)

type Connection struct {
	User1Id   int64     `db:"user1_id"`
	User2Id   int64     `db:"user2_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (db *MyDb) newConnection(user1Id, user2Id string) (*Connection, error) {
	id1, id2, err := db.validateConnectionId(user1Id, user2Id)
	if err != nil {
		return nil, err
	}
	return &Connection{id1, id2, time.Now()}, err
}

func (db *MyDb) PostConnection(conn *Connection) error {
	err := db.Insert(conn)
	return err
}

func (db *MyDb) DeleteConnection(user1Id, user2Id string) error {
	//  id1, id2, err := db.validateConnectionId(user1Id, user2Id)
	// TODO
	return nil
}
