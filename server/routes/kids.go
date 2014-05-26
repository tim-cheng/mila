package routes

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"strconv"
	"time"
)

func (rt *Routes) PostKid(params martini.Params, r render.Render, req *http.Request) {
	bdForm := "2006-01-02"
	bd, _ := time.Parse(bdForm, req.FormValue("birthday"))

	kid, err := rt.Db.NewKid(
		params["id"],
		req.FormValue("name"),
		bd,
		(req.FormValue("type") == "boy"),
	)

	if err == nil {
		err = rt.Db.PostKid(kid)
	}

	if err == nil {
		go rt.UpdateUserDesc(params["id"])
		r.JSON(201, map[string]interface{}{
			"id": kid.Id,
		})
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to add kid " + err.Error(),
		})
	}
}

func (rt *Routes) GetKids(params martini.Params, r render.Render) {
	kids, err := rt.Db.GetKids(params["id"])
	if err == nil && len(kids) > 0 {
		retKids := make([]map[string]interface{}, len(kids))
		for i := range kids {
			retKids[i] = map[string]interface{}{
				"id":   kids[i].Id,
				"name": kids[i].Name,
				"age":  birthdayToAge(kids[i].Birthday),
				"boy":  kids[i].IsBoy,
			}
		}
		r.JSON(200, retKids)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to get kids ",
		})
	}
}

func (rt *Routes) DeleteKid(params martini.Params, r render.Render) {
	for {
		kid, err := rt.Db.GetKid(params["kid"])
		if err != nil {
			break
		}
		pId, err := strconv.Atoi(params["id"])
		if err != nil || kid.ParentId != int64(pId) {
			break
		}
		err = rt.Db.DeleteKid(kid)
		if err != nil {
			break
		}
		go rt.UpdateUserDesc(params["id"])
		r.JSON(200, nil)
		return
	}
	r.JSON(404, map[string]interface{}{
		"message": "Failed to delete kid",
	})
}

func birthdayToAge(birthday time.Time) int {
	age := time.Now().Sub(birthday).Hours() / (24 * 365)
	return int(age)
}
