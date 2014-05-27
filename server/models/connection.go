package models

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"time"
)

type Connection struct {
	User1Id   int64     `db:"user1_id"`
	User2Id   int64     `db:"user2_id"`
	CreatedAt time.Time `db:"created_at"`
}

// Validation Hooks
func (c *Connection) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) NewConnection(user1Id, user2Id string) (*Connection, error) {
	id1, id2, err := db.validateConnectionId(user1Id, user2Id)
	if err != nil {
		return nil, err
	}
	return &Connection{User1Id: id1, User2Id: id2}, err
}

func (db *MyDb) update2dConnection(userId int64) {
	ids, err := db.Get1dConnectionById(userId)
	if err == nil {
		set2d := make(map[int64]bool)
		for _, id := range ids {
			id2s, err := db.Get1dConnectionById(id)
			if err == nil {
				for _, v := range id2s {
					set2d[v] = true
				}
			}
		}
		delete(set2d, userId)
		db.Update2dConnection(userId, len(set2d))
	}
}

func (db *MyDb) updateConnections(user1Id, user2Id int64) {
	// 1degree:
	nConn, err := db.GetNumConnections(user1Id)
	if err == nil {
		db.Update1dConnection(user1Id, nConn)
	}
	nConn, err = db.GetNumConnections(user2Id)
	if err == nil {
		db.Update1dConnection(user2Id, nConn)
	}

	// 2degree:
	// find all 1d connection, which will need to update 2d stat
	set2d := make(map[int64]bool)
	ids, err := db.Get1dConnectionById(user1Id)
	if err == nil {
		for _, v := range ids {
			set2d[v] = true
		}
	}
	ids, err = db.Get1dConnectionById(user2Id)
	if err == nil {
		for _, v := range ids {
			set2d[v] = true
		}
	}
	// update 2d connection for every node
	for id, _ := range set2d {
		db.update2dConnection(id)
	}
}


func (db *MyDb) PostConnection(conn *Connection) error {
	err := db.Insert(conn)
	if err == nil {
		go db.updateConnections(conn.User1Id, conn.User2Id)
	}
	return err
}

func (db *MyDb) DeleteConnection(user1Id, user2Id string) error {
	//  id1, id2, err := db.validateConnectionId(user1Id, user2Id)
	// TODO
	return nil
}

func (db *MyDb) GetNumConnections(userId int64) (int, error) {
	count, err := db.SelectInt("select count(*) from connections where (user1_id=$1 or user2_id=$1)", userId)
	return int(count), err
}

func (db *MyDb) GetConnections(userId string) ([]interface{}, error) {
	uId, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	return db.Select(Connection{}, "select user1_id, user2_id from connections where user1_id=$1 or user2_id=$1", uId)
}

func (db *MyDb) Get1dConnectionById(userId int64) ([]int64, error) {
	var userIds []int64
	_, err := db.Select(&userIds, "select user1_id from connections where user2_id=$1 union select user2_id from connections where user1_id=$1", userId)
	fmt.Println("Get1dConnectionById ", userIds, err)
	return userIds, err
}
