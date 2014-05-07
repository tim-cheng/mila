package routes

import (
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
)

func (rt *Routes) PostConnection(r render.Render, req *http.Request) {
	conn, err := rt.Db.NewConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
	if err == nil {
		err = rt.Db.PostConnection(conn)
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
}

func (rt *Routes) DeleteConnection(r render.Render, req *http.Request) {
	err := rt.Db.DeleteConnection(req.FormValue("user1_id"), req.FormValue("user2_id"))
	if err == nil {
		r.JSON(200, nil)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to delete connection " + err.Error(),
		})
	}
}
