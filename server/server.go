package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	//"fmt"
)

var myDb *MyDb

func main() {
	myDb = NewDb()
	defer myDb.Db.Close()

	m := martini.Classic()
	m.Use(render.Renderer())

	// Routes

	m.Get("/", func() string {
		return "Welcome to Mila"
	})

	// users
	m.Get("/users/:id", func(params martini.Params, r render.Render) {
		user, err := myDb.GetUser(params["id"])
		renderResponse(r, err, 200, user, 404, "User not found")
	})

	m.Post("/users", func(r render.Render, req *http.Request) {
		user, err := myDb.newUser("email", req.FormValue("email"), req.FormValue("first_name"), req.FormValue("last_name"))
		if err == nil {
			err = myDb.PostUser(user)
		}
		renderResponse(r, err, 201, user, 404, "Failed to add user")
	})

	// connections
	m.Post("/connections", func(r render.Render, req *http.Request) {
		conn, err := myDb.newConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
		if err == nil {
			err = myDb.PostConnection(conn)
		}
		renderResponse(r, err, 201, conn, 404, "Failed to add connection")
	})

	m.Delete("/connections", func(r render.Render, req *http.Request) {
		err := myDb.DeleteConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
		renderResponse(r, err, 200, nil, 404, "Failed to add connection")
	})

	http.ListenAndServe(":8080", m)
}

func renderResponse(r render.Render, err error, passCode int, passObj interface{}, failCode int, failMsg string) {
	if err == nil {
		r.JSON(passCode, passObj)
	} else {
		r.JSON(404, map[string]interface{}{"message": failMsg})
	}
}
