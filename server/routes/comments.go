package routes

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/tim-cheng/mila/server/models"
	"net/http"
)

func (rt *Routes) PostComment(params martini.Params, r render.Render, req *http.Request) {
	c, err := rt.Db.NewComemnt(
		req.FormValue("user_id"),
		params["id"],
		req.FormValue("body"),
	)
	if err == nil {
		err = rt.Db.PostComment(c)
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
}

func (rt *Routes) GetComments(params martini.Params, r render.Render) {
	comments, err := rt.Db.GetComments(params["id"])
	if err == nil && len(comments) > 0 {
		retComments := make([]map[string]interface{}, len(comments))
		for i := range comments {
			c := comments[i].(*models.Comment)
			retComments[i] = map[string]interface{}{
				"id":         c.Id,
				"user_id":    c.UserId,
				"post_id":    c.PostId,
				"body":       c.Body,
				"created_at": c.CreatedAt,
			}
		}
		r.JSON(200, retComments)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to get comments",
		})
	}
}
