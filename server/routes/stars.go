package routes

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
)

func (rt *Routes) PutStar(params martini.Params, r render.Render, req *http.Request) {
	s, err := rt.Db.NewStar(
		req.FormValue("user_id"),
		params["id"],
	)
	if err == nil {
		err = rt.Db.PutStar(s)
	}
	if err == nil {
		go func() {
			p, err := rt.Db.GetPost(params["id"])
			if err == nil {
				if p.UserId != s.UserId {
					rt.postActivityLike(s.UserId, p.UserId, p.Id)
				}
			}
		}()
		r.JSON(200, nil)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to add star " + err.Error(),
		})
	}
}

func (rt *Routes) DeleteStar(params martini.Params, r render.Render, req *http.Request) {
	s, err := rt.Db.NewStar(
		req.FormValue("user_id"),
		params["id"],
	)
	if err == nil {
		err = rt.Db.DeleteStar(s)
	}
	if err == nil {
		r.JSON(200, nil)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to delete star " + err.Error(),
		})
	}
}
