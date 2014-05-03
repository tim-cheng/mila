package main

import (
	"errors"
	"strconv"
)

func (db *MyDb) validateUserId(userId string) (id int64, err error) {
	// validate ids are valid numbers
	id, err = strconv.ParseInt(userId, 0, 64)
	if err != nil {
		return
	}
	// validate user exists
	_, err = db.Get(User{}, id)
	if err != nil {
		return
	}
	return
}

func (db *MyDb) validatePostId(postId string) (id int64, err error) {
	// validate ids are valid numbers
	id, err = strconv.ParseInt(postId, 0, 64)
	if err != nil {
		return
	}
	// validate user exists
	_, err = db.Get(Post{}, id)
	if err != nil {
		return
	}
	return
}

func (db *MyDb) validateConnectionId(user1Id string, user2Id string) (id1 int64, id2 int64, err error) {

	// validate ids are valid numbers
	id1, err = db.validateUserId(user1Id)
	if err != nil {
		return
	}
	id2, err = db.validateUserId(user2Id)
	if err != nil {
		return
	}

	// validate id1/id2
	if id1 == id2 {
		err = errors.New("connection ids same")
	}
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	return
}
