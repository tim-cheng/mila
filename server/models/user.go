package models

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"strconv"
	"time"
)

type User struct {
	Id          int64     `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Type        string    `db:"type"`
	FbId        string    `db:"fb_id"`
	Admin       bool      `db:"admin"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	NumDegree1  int64     `db:"num_degree1"`
	NumDegree2  int64     `db:"num_degree2"`
	Description string    `db:"description"`
	Picture     []byte    `db:"picture"`
}

// Validation Hooks
func (u *User) PreInsert(s gorp.SqlExecutor) error {
	u.CreatedAt = time.Now()
	return nil
}

func (db *MyDb) NewUser(typ, email, password, firstName, lastName, fb_id string) (*User, error) {

	// generate hash
	hashQuery := fmt.Sprintf("select crypt('%s', gen_salt('md5'))", password)
	var hash string
	db.SelectOne(&hash, hashQuery)

	return &User{
		Type:        typ,
		Email:       email,
		Password:    hash,
		Admin:       false,
		FirstName:   firstName,
		LastName:    lastName,
		Description: "proud parent",
		FbId:        fb_id,
	}, nil
}

func (db *MyDb) GetUser(userId string) (*User, error) {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}
	u := new(User)
	err = db.SelectOne(u, "select id, created_at, email, type, first_name, last_name, num_degree1, num_degree2, description from users where id=$1", id)
	return u, err
}

func (db *MyDb) GetUserByEmail(email string) (*User, error) {
	u := new(User)
	err := db.SelectOne(u, "select id, created_at, email, type, first_name, last_name, num_degree1, num_degree2, description from users where email=$1", email)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (db *MyDb) PostUser(user *User) error {
	err := db.Insert(user)
	return err
}

func (db *MyDb) PostUserPicture(userId int64, image []byte) error {
	u := new(User)
	err := db.SelectOne(u, "select * from users where id=$1", userId)
	u.Picture = image
	_, err = db.Update(u)
	return err
}

func (db *MyDb) GetUserPicture(userId int64) ([]byte, error) {
	u := new(User)
	err := db.SelectOne(u, "select picture from users where id=$1", userId)
	return u.Picture, err
}

func (db *MyDb) GetPassword(email string) (string, error) {
	p := ""
	err := db.SelectOne(&p, "select password from users where email=$1", email)
	return p, err
}
