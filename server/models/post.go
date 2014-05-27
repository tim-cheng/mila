package models

import (
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"time"
)

type Post struct {
	Id         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UserId     int64     `db:"user_id"`
	Body       string    `db:"body"`
	BgColor    string    `db:"bg_color"`
	Picture    []byte    `db:"picture"`
	HasPicture bool      `db:"has_picture"`
}

type PostFeed struct {
	Id         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UserId     int64     `db:"user_id"`
	Body       string    `db:"body"`
	BgColor    string    `db:"bg_color"`
	Picture    []byte    `db:"picture"`
	HasPicture bool      `db:"has_picture"`
	RefUserId  int64     `db:"ref_user_id"`
}

// Validation Hooks
func (p *Post) PreInsert(s gorp.SqlExecutor) error {
	p.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) NewPost(userId, content, bgcolor string) (*Post, error) {
	id, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	return &Post{
		UserId:  id,
		Body:    content,
		BgColor: bgcolor,
	}, nil
}

func (db *MyDb) PostPost(post *Post) error {
	err := db.Insert(post)
	return err
}

func (db *MyDb) GetPost(postId string) (*Post, error) {
	pId, err := db.validatePostId(postId)
	if err != nil {
		return nil, err
	}
	p := new(Post)
	err = db.SelectOne(p, "select id, user_id from posts where id=$1", pId)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (db *MyDb) DeletePost(postId string) error {
	p, err := db.GetPost(postId)
	if err != nil {
		return err
	}
	// TODO: should this be done in the background?
	_, err = db.Exec("delete from comments where post_id=$1", p.Id)
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from stars where post_id=$1", p.Id)
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from activities where post_id=$1", p.Id)
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from feeds where post_id=$1", p.Id)
	if err != nil {
		return err
	}
	count, err := db.Delete(p)
	if count != 1 {
		return errors.New("couldn't delete post")
	}
	return err
}

func (db *MyDb) GetPosts(userId string, degree string) ([]interface{}, error) {
	id, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}

	var posts []interface{}
	if degree == "" || degree == "0" {
		posts, err = db.Select(PostFeed{}, "select id Id, created_at CreatedAt, user_id UserId, body Body, bg_color BgColor, has_picture HasPicture from posts where user_id=$1", id)
		fmt.Println("!!!!!! errorrr ", posts, err)
	} else if degree == "1" {
		posts, err = db.Select(PostFeed{},
			"select id Id, created_at CreatedAt, user_id UserId, body Body, bg_color BgColor, has_picture HasPicture from posts where user_id in "+
				"(select $1 UNION (select user2_id from connections where user1_id=$1) "+
				"UNION (select user1_id from connections where user2_id=$1)) "+
				"order by created_at desc", id)
		fmt.Println("!!!!!! errorrr ", posts, err)
	} else if degree == "2" {
		posts, err = db.Select(PostFeed{},
			"select posts.id Id, posts.created_at CreatedAt, posts.user_id UserId, posts.body Body, posts.bg_color BgColor, posts.has_picture HasPicture, feeds.ref_user_id RefUserId from posts "+
				"join feeds on posts.id = feeds.post_id "+
				"where feeds.user_id=$1", id)
	}

	return posts, err
}

func (db *MyDb) PostPostPicture(postId string, image []byte) error {
	pId, err := db.validatePostId(postId)
	if err != nil {
		return err
	}
	p := new(Post)
	err = db.SelectOne(p, "select id, created_at, user_id, bg_color, body from posts where id=$1", pId)
	p.Picture = image
	p.HasPicture = true
	_, err = db.Update(p)
	return err
}

func (db *MyDb) GetPostPicture(postId string) ([]byte, error) {
	pId, err := db.validatePostId(postId)
	if err != nil {
		return nil, err
	}
	p := new(Post)
	err = db.SelectOne(p, "select picture from posts where id=$1", pId)
	return p.Picture, err
}
