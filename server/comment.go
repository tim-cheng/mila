package main

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Comment struct {
	Id        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	PostId    int64     `db:"post_id"`
	UserId    int64     `db:"user_id"`
	Body      string    `db:"body"`
}

// Validation Hooks
func (c *Comment) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) newComemnt(userId, postId, content string) (*Comment, error) {
	uId, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	pId, err := db.validatePostId(postId)
	if err != nil {
		return nil, err
	}
	return &Comment{
		PostId: pId,
		UserId: uId,
		Body:   content,
	}, nil
}

func (db *MyDb) PostComment(c *Comment) error {
	err := db.Insert(c)
	return err
}

func (db *MyDb) GetNumComments(postId int64) (int, error) {
	count, err := db.SelectInt("select count(*) from comments where post_id=$1", postId)
	return int(count), err
}
