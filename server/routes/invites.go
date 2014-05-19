package routes

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"strings"
)

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
