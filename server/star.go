package main

import (
	//  "fmt"
	"errors"
)

type Star struct {
	PostId int64 `db:"post_id"`
	UserId int64 `db:"user_id"`
}

func (db *MyDb) newStar(userId, postId string) (*Star, error) {
	uId, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	pId, err := db.validatePostId(postId)
	if err != nil {
		return nil, err
	}
	return &Star{
		PostId: pId,
		UserId: uId,
	}, nil
}

func (db *MyDb) PutStar(s *Star) error {
	err := db.Insert(s)
	return err
}

func (db *MyDb) DeleteStar(s *Star) error {
	count, err := db.Delete(s)
	if count != 1 {
		return errors.New("couldn't delete star")
	}
	return err
}

func (db *MyDb) GetNumStars(postId int64) (int, error) {
	count, err := db.SelectInt("select count(*) from stars where post_id=$1", postId)
	return int(count), err
}
