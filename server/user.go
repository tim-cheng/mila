package main

import (
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Id          int64     `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Type        string    `db:"type"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	NumDegree1  int64     `db:"num_degree1"`
	NumDegree2  int64     `db:"num_degree2"`
	Description string    `db:"description"`
	PictureUrl  string    `db:"picture_url"`
}

func (db *MyDb) newUser(typ, email, password, firstName, lastName string) (*User, error) {

	// generate hash
	hashQuery := fmt.Sprintf("select crypt('%s', gen_salt('md5'))", password)
	var hash string
	db.SelectOne(&hash, hashQuery)
	fmt.Println("hash = ", hash)

	return &User{
		Type:      typ,
		Email:     email,
		Password:  hash,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
	}, nil
}

func (db *MyDb) GetUser(userId string) (*User, error) {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}
	u := new(User)
	err = db.SelectOne(u, "select * from users where id=$1", id)
	return u, err
}

func (db *MyDb) PostUser(user *User) error {
	err := db.Insert(user)
	return err
}

func (db *MyDb) GetPassword(email string) (string, error) {
	p := ""
	err := db.SelectOne(&p, "select password from users where email=$1", email)
	return p, err
}
