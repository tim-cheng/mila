package models

import (
	"errors"
	"github.com/coopernurse/gorp"
	"strconv"
	"time"
)

type Kid struct {
	Id        int64     `db:"id"`
	ParentId  int64     `db:"parent_id"`
	Name      string    `db:"name"`
	Birthday  time.Time `db:"birthday"`
	IsBoy     bool      `db:"is_boy"`
	CreatedAt time.Time `db:"created_at"`
	Picture   []byte    `db:"picture"`
}

func (k *Kid) PreInsert(s gorp.SqlExecutor) error {
	k.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) NewKid(parentId string, name string, birthday time.Time, isBoy bool) (*Kid, error) {
	pId, err := db.validateUserId(parentId)
	if err != nil {
		return nil, err
	}
	return &Kid{
		ParentId: pId,
		Name:     name,
		Birthday: birthday,
		IsBoy:    isBoy,
	}, nil
}

func (db *MyDb) PostKid(kid *Kid) error {
	err := db.Insert(kid)
	return err
}

func (db *MyDb) GetKids(parentId string) ([]Kid, error) {
	uId, err := db.validateUserId(parentId)
	if err != nil {
		return nil, err
	}
	var kids []Kid
	_, err = db.Select(&kids, "select id, parent_id, name, birthday, is_boy from kids where parent_id=$1", uId)
	return kids, err
}

func (db *MyDb) GetKid(kidId string) (*Kid, error) {
	id, err := strconv.ParseInt(kidId, 0, 64)
	if err != nil {
		return nil, err
	}
	var kid Kid
	err = db.SelectOne(&kid, "select id, parent_id from kids where id=$1", id)
	return &kid, err
}

func (db *MyDb) DeleteKid(k *Kid) error {
	count, err := db.Delete(k)
	if count != 1 {
		return errors.New("couldn't delete post")
	}
	return err
}
