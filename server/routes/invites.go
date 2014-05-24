package routes

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/tim-cheng/mila/server/models"
	"strings"
)

func (rt *Routes) GetInvites(params martini.Params, r render.Render) {
	invites, err := rt.Db.GetInvites(params["id"])
	if err == nil && len(invites) > 0 {
		retInvites := make([]map[string]interface{}, len(invites))
		for i := range invites {
			inv := invites[i].(*models.Invite)
			u, _ := rt.Db.GetUserName(inv.User1Id)
			retInvites[i] = map[string]interface{}{
				"user_id":    inv.User1Id,
				"create_at":  inv.CreatedAt,
				"first_name": u.FirstName,
				"last_name":  u.LastName,
			}
		}
		r.JSON(200, retInvites)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to get invites ",
		})
	}
}

func (rt *Routes) DeleteInvite(params martini.Params, r render.Render) {
	for {
		inv, err := rt.Db.NewInvite(params["id"], params["id2"])
		if err != nil {
			break
		}
		err = rt.Db.DeleteInvite(inv)
		if err != nil {
			break
		}

		conn, err := rt.Db.NewConnection(params["id"], params["id2"])
		if err != nil {
			break
		}

		err = rt.Db.PostConnection(conn)
		if err != nil {
			break
		}
		r.JSON(200, map[string]interface{}{
			"message": "connected",
		})
		return
	}
	r.JSON(404, map[string]interface{}{
		"message": "failed to accept invite",
	})
}

func (rt *Routes) PostInvite(params martini.Params, r render.Render) {
	for {
		inv, err := rt.Db.NewInvite(params["id"], params["id2"])
		if err != nil {
			break
		}

		err = rt.Db.PostInvite(inv)
		if err != nil {
			break
		}

		u1, err := rt.Db.GetUser(params["id"])
		if err != nil {
			break
		}

		u2, err := rt.Db.GetUser(params["id2"])
		if err != nil {
			break
		}

		user2Email := u2.Email
		user2Email = strings.TrimSuffix(user2Email, "@fb")
		go sendUserInviteEmail(user2Email, u2.FirstName, u1.FirstName+" "+u1.LastName)
		go sendUserPushMsg(u2.Id, u1.FirstName+" "+u1.LastName+" would like to connect")
		go rt.postActivityInvite(u1.Id, u2.Id, "hello there!")

		r.JSON(201, map[string]interface{}{
			"user1_id": inv.User1Id,
			"user2_id": inv.User2Id,
		})
		return
	}

	r.JSON(404, map[string]interface{}{
		"message": "Failed to add invite",
	})
}
