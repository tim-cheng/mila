package models

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Activity struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	FriendId  int64     `db:"friend_id"`
	Type      int       `db:"type"`
	Message   string    `db:"message"`
	PostId    int64     `db:"post_id"`
	CreatedAt time.Time `db:"created_at"`
}

const (
	ActivityTypePost    = 1
	ActivityTypeComment = 2
	ActivityTypeLike    = 3
	ActivityTypeInvite  = 4
)

// Validation Hooks
func (a *Activity) PreInsert(s gorp.SqlExecutor) error {
	a.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) PostActivity(a *Activity) error {
	err := db.Insert(a)
	return err
}

func (db *MyDb) GetActivities(userId string) ([]interface{}, error) {
	id, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	activities, err := db.Select(Activity{}, "select * from activities where user_id=$1 order by id desc", id)
	return activities, err
}

func (db *MyDb) NewActivity(userId int64, friendId int64, postId int64, typ int, msg string) (*Activity, error) {
	return &Activity{
		UserId:   userId,
		FriendId: friendId,
		Type:     typ,
		Message:  msg,
		PostId:   postId,
	}, nil
}
