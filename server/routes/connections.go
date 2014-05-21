package routes

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/tim-cheng/mila/server/models"
	"net/http"
	"strconv"
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

func (rt *Routes) GetConnections(params martini.Params, r render.Render) {
	userId := params["id"]
	uId, err := strconv.ParseInt(userId, 0, 64)
	if err != nil {
		r.JSON(404, map[string]interface{}{
			"message": "invalid user_id",
		})
	}

	conns, err := rt.Db.GetConnections(userId)
	if err == nil && len(conns) > 0 {
		retConns := make([]map[string]interface{}, len(conns))
		for i := range conns {
			conn := conns[i].(*models.Connection)
			if uId == conn.User1Id {
				uId = conn.User2Id
			} else {
				uId = conn.User1Id
			}
			u, _ := rt.Db.GetUserName(uId)
			retConns[i] = map[string]interface{}{
				"user_id":    uId,
				"create_at":  conn.CreatedAt,
				"first_name": u.FirstName,
				"last_name":  u.LastName,
			}
		}
		r.JSON(200, retConns)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to get connections ",
		})
	}
}
