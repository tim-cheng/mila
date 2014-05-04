package main

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Post struct {
	Id         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UserId     int64     `db:"user_id"`
	Body       string    `db:"body"`
	PictureUrl string    `db:"picture_url"`
}

// Validation Hooks
func (p *Post) PreInsert(s gorp.SqlExecutor) error {
	p.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) newPost(userId, content string) (*Post, error) {
	id, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	return &Post{
		UserId: id,
		Body:   content,
	}, nil
}

func (db *MyDb) PostPost(post *Post) error {
	err := db.Insert(post)
	return err
}

func (db *MyDb) GetPosts(userId string) ([]interface{}, error) {
	id, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	posts, err := db.Select(Post{}, "select * from posts where user_id=$1", id)
	return posts, err
}
