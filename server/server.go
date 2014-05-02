package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"strconv"
)

func main() {
	dbmap := initDb()
	defer dbmap.Db.Close()

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func() string {
		return "Welcome to Mila"
	})

	m.Get("/users/:id", func(params martini.Params, r render.Render) {
		paramId, err := strconv.Atoi(params["id"])
		if err != nil {
			r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": err.Error()})
			return
		}
		user, error := GetUser(dbmap, paramId)
		if error == nil {
			r.JSON(200, map[string]interface{}{"status": "Success", "data": user})
		} else {
			r.JSON(404, map[string]interface{}{"status": "Fail", "error_message": "User not found"})
		}
	})

	http.ListenAndServe(":8080", m)
}
