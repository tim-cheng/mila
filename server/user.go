package main

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"time"
)

type User struct {
	Id          int64     `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Email       string    `db:"email"`
	Type        string    `db:"type"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	NumDegree1  int64     `db:"num_degree1"`
	NumDegree2  int64     `db:"num_degree2"`
	Description string    `db:"description"`
	PictureUrl  string    `db:"picture_url"`
}

func newUser(typ, email, firstName, lastName string) User {
	return User{
		Type:      typ,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
	}
}

func GetUser(dbmap *gorp.DbMap, id int) (*User, error) {
	u := new(User)
	err := dbmap.SelectOne(u, "select * from users where id=$1", id)
	fmt.Println("user = ", u, " err= ", err)
	return u, err
}
