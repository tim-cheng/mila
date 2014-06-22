package main

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/tim-cheng/mila/server/models"
	"github.com/tim-cheng/mila/server/routes"
	"net/http"
)

var myDb *models.MyDb

func Secret(user, realm string) string {
	p, err := myDb.GetPassword(user)
	if err != nil {
		return ""
	} else {
		return p
	}
}

func startServer() {
	myDb = models.NewDb()
	defer myDb.Db.Close()

	m := martini.Classic()
	m.Use(render.Renderer())

	router := routes.New(myDb)

	// authentication
	basicAuth := auth.NewBasicAuthenticator("mila.com", Secret)
	authFunc := basicAuth.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		fmt.Println("auth user: ", r.Username)
	})
	//m.Use(authFunc)

	// Routes
	m.Use(martini.Static("assets"))

	m.Get("/login", authFunc, router.Login)
	m.Get("/login_facebook", router.LoginFacebook)

	m.Get("/users/:id", authFunc, router.GetUser)
	m.Put("/users/:id", authFunc, router.PutUser)
	m.Post("/users/:id/picture", authFunc, router.PostUserPicture)
	// TODO: get image doesn't require basic auth to make it easier to fetch/cache
	m.Get("/users/:id/picture", router.GetUserPicture)
	m.Post("/users", router.PostUser)

	m.Get("/users", authFunc, router.SearchUsers)

	m.Get("/connections/:id", authFunc, router.GetConnections)
	m.Post("/connections", authFunc, router.PostConnection)
	m.Delete("/connections", authFunc, router.DeleteConnection)

	m.Get("/users/:id/invite", authFunc, router.GetInvites)
	m.Post("/users/:id/invite/:id2", authFunc, router.PostInvite)
	m.Delete("/users/:id/invite/:id2", authFunc, router.DeleteInvite)

	m.Post("/users/:id/fb_invite", authFunc, router.PostFbInvite)

	m.Post("/posts/:id/picture", authFunc, router.PostPostPicture)
	m.Get("/posts/:id/picture", authFunc, router.GetPostPicture)

	m.Post("/posts", authFunc, router.PostPost)
	m.Get("/posts", authFunc, router.GetPosts)
	m.Delete("/posts/:id", authFunc, router.DeletePost)

	m.Post("/posts/:id/comments", authFunc, router.PostComment)
	m.Get("/posts/:id/comments", authFunc, router.GetComments)

	m.Put("/posts/:id/stars", authFunc, router.PutStar)
	m.Delete("/posts/:id/stars", authFunc, router.DeleteStar)

	m.Get("/activities/:id", authFunc, router.GetActivities)

	m.Post("/users/:id/kids", router.PostKid)
	m.Get("/users/:id/kids", router.GetKids)
	m.Delete("/users/:id/kids/:kid", router.DeleteKid)

	http.ListenAndServe(":8080", m)
}

func main() {
	go startServer()
	select {}
}
