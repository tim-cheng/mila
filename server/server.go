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

func startServer() {
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

	// login
	m.Get("/login", authFunc, func(r render.Render, req *http.Request) {
		email := basicAuth.CheckAuth(req)
		user, err := myDb.GetUserByEmail(email)
		fmt.Println("email: ", email, " user.id = ", user.Id)
		if err == nil {
			r.JSON(200, map[string]interface{}{
				"id": user.Id,
			})
		} else {
			r.JSON(500, nil)
		}
	})

	// users
	m.Get("/users/:id", authFunc, func(params martini.Params, r render.Render) {
		user, err := myDb.GetUser(params["id"])
		if err == nil {
			r.JSON(200, map[string]interface{}{
				"id":          user.Id,
				"first_name":  user.FirstName,
				"last_name":   user.LastName,
				"email":       user.Email,
				"description": user.Description,
				"num_degree1": user.NumDegree1,
				"num_degree2": user.NumDegree2,
			})
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "User not found " + err.Error(),
			})
		}
	})

	m.Post("/users", func(r render.Render, req *http.Request) {
		user, err := myDb.newUser(
			"email",
			req.FormValue("email"),
			req.FormValue("password"),
			req.FormValue("first_name"),
			req.FormValue("last_name"),
		)

		if err == nil {
			err = myDb.PostUser(user)
		}
		if err == nil {
			r.JSON(201, map[string]interface{}{
				"id": user.Id,
			})
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to add user " + err.Error(),
			})
		}
	})

	// connections
	m.Post("/connections", authFunc, func(r render.Render, req *http.Request) {
		conn, err := myDb.newConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
		if err == nil {
			err = myDb.PostConnection(conn)
		}
		if err == nil {
			r.JSON(201, map[string]interface{}{
				"user1_id": conn.User1Id,
				"user2_id": conn.User2Id,
			})
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to add connection " + err.Error(),
			})
		}
	})

	m.Delete("/connections", authFunc, func(r render.Render, req *http.Request) {
		err := myDb.DeleteConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
		if err == nil {
			r.JSON(200, nil)
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to delete connection " + err.Error(),
			})
		}
	})

	// posts
	m.Post("/posts", authFunc, func(r render.Render, req *http.Request) {
		post, err := myDb.newPost(req.FormValue("user_id"), req.FormValue("body"))
		if err == nil {
			err = myDb.PostPost(post)
		}
		if err == nil {
			r.JSON(201, map[string]interface{}{
				"id": post.Id,
			})
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to add post " + err.Error(),
			})
		}
	})

	m.Get("/posts", authFunc, func(r render.Render, req *http.Request) {
		posts, err := myDb.GetPosts(req.FormValue("user_id"))
		if err == nil {
			r.JSON(200, posts)
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to get posts " + err.Error(),
			})
		}
	})

	// comments
	m.Post("/posts/:id/comments", authFunc, func(params martini.Params, r render.Render, req *http.Request) {
		c, err := myDb.newComemnt(
			req.FormValue("user_id"),
			params["id"],
			req.FormValue("body"),
		)
		if err == nil {
			err = myDb.PostComment(c)
		}
		if err == nil {
			r.JSON(201, map[string]interface{}{
				"id": c.Id,
			})
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to add comment " + err.Error(),
			})
		}
	})

	// stars
	m.Put("/posts/:id/stars", authFunc, func(params martini.Params, r render.Render, req *http.Request) {
		s, err := myDb.newStar(
			req.FormValue("user_id"),
			params["id"],
		)
		if err == nil {
			err = myDb.PutStar(s)
		}
		if err == nil {
			r.JSON(200, nil)
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to add star " + err.Error(),
			})
		}
	})

	m.Delete("/posts/:id/stars", authFunc, func(params martini.Params, r render.Render, req *http.Request) {
		s, err := myDb.newStar(
			req.FormValue("user_id"),
			params["id"],
		)
		if err == nil {
			err = myDb.DeleteStar(s)
		}
		if err == nil {
			r.JSON(200, nil)
		} else {
			r.JSON(404, map[string]interface{}{
				"message": "Failed to delete star " + err.Error(),
			})
		}
	})

	http.ListenAndServe(":8080", m)
}

func main() {
	go startServer()
	select {}
}
