package models

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Invite struct {
	User1Id   int64     `db:"user1_id"`
	User2Id   int64     `db:"user2_id"`
	CreatedAt time.Time `db:"created_at"`
}

// Validation Hooks
func (inv *Invite) PreInsert(s gorp.SqlExecutor) error {
	inv.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) NewInvite(user1Id, user2Id string) (*Invite, error) {
	id1, id2, err := db.validateConnectionId(user1Id, user2Id)
	if err != nil {
		return nil, err
	}
	return &Invite{User1Id: id1, User2Id: id2}, err
}

func (db *MyDb) PostInvite(inv *Invite) error {
	// TODO: ignore post error....
	db.Insert(inv)
	return nil
}

func (db *MyDb) GetInvites(userId string) ([]interface{}, error) {
	uId, err := db.validateUserId(userId)
	if err != nil {
		return nil, err
	}
	return db.Select(Invite{}, "select user1_id, user2_id from invites where user2_id=$1", uId)
}

func (db *MyDb) DeleteInvite(user1Id, user2Id string) error {
	//  id1, id2, err := db.validateConnectionId(user1Id, user2Id)
	// TODO
	return nil
}
