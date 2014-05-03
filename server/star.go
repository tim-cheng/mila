package main

/*import (
  "fmt"
)*/

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
