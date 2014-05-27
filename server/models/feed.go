package models

import (
//"github.com/coopernurse/gorp"
)

type Feed struct {
	UserId    int64 `db:"user_id"`
	PostId    int64 `db:"post_id"`
	RefUserId int64 `db:"ref_user_id"`
}

func (db *MyDb) PostFeed(f *Feed) error {
	err := db.Insert(f)
	return err
}

func (db *MyDb) New1dFeed(userId, postId int64) (*Feed, error) {
	return &Feed{
		UserId:    userId,
		PostId:    postId,
		RefUserId: userId,
	}, nil
}

func (db *MyDb) New2dFeed(userId, postId, refUserId int64) (*Feed, error) {
	return &Feed{
		UserId:    userId,
		PostId:    postId,
		RefUserId: refUserId,
	}, nil
}
