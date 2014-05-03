package main

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
)

var myDb *MyDb

func Secret(user, realm string) string {
	p, err := myDb.GetPassword(user)
	if err != nil {
		return ""
	} else {
		return p
	}
}

// injection example (injecting username), not used
type authUser struct {
	User string
}
type AuthUser interface {
	GetUser() string
}

func (au *authUser) GetUser() string {
	return au.User
}

func main() {
	myDb = NewDb()
	defer myDb.Db.Close()

	m := martini.Classic()
	m.Use(render.Renderer())

	// authentication
	basicAuth := auth.NewBasicAuthenticator("mila.com", Secret)
	authFunc := basicAuth.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		fmt.Println("auth user: ", r.Username)
	})
	//m.Use(authFunc)

	// request level injection example
	m.Use(func(req *http.Request, c martini.Context) {
		// inject to interface (override existing interface)
		//c.MapTo(&authUser{basicAuth.CheckAuth(req)}, (*AuthUser)(nil))
		// inject to struct
		c.Map(&authUser{basicAuth.CheckAuth(req)})
	})

	// Routes
	m.Get("/", func() string {
		return "Welcome to Mila"
	})

	// users
	m.Get("/users/:id", authFunc, func(params martini.Params, r render.Render) {
		user, err := myDb.GetUser(params["id"])
		renderResponse(r, err, 200, user, 404, "User not found")
	})

	m.Post("/users", func(r render.Render, req *http.Request) {
		user, err := myDb.newUser("email", req.FormValue("email"), req.FormValue("password"), req.FormValue("first_name"), req.FormValue("last_name"))
		if err == nil {
			err = myDb.PostUser(user)
		}
		renderResponse(r, err, 201, user, 404, "Failed to add user")
	})

	// connections
	m.Post("/connections", authFunc, func(r render.Render, req *http.Request) {
		conn, err := myDb.newConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
		if err == nil {
			err = myDb.PostConnection(conn)
		}
		renderResponse(r, err, 201, conn, 404, "Failed to add connection")
	})

	m.Delete("/connections", authFunc, func(r render.Render, req *http.Request) {
		err := myDb.DeleteConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
		renderResponse(r, err, 200, nil, 404, "Failed to add connection")
	})

	// posts
	m.Post("/posts", authFunc, func(r render.Render, req *http.Request) {
		post, err := myDb.newPost(req.FormValue("user_id"), req.FormValue("body"))
		if err == nil {
			err = myDb.PostPost(post)
		}
		renderResponse(r, err, 201, post, 404, "Failed to add post")
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
