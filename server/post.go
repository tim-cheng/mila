package main

import (
	"time"
)

type Post struct {
	Id         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UserId     int64     `db:"user_id"`
	Body       string    `db:"body"`
	PictureUrl string    `db:"picture_url"`
}

func (db *MyDb) newPost(userId, content string) (*Post, error) {
	id, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	return &Post{
		CreatedAt: time.Now(),
		UserId:    id,
		Body:      content,
	}, nil
}

func (db *MyDb) PostPost(post *Post) error {
	err := db.Insert(post)
	return err
}
