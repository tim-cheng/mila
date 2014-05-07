package routes

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// from https://codereview.appspot.com/76540043/patch/80001/90001
func basicAuth(r *http.Request) (username, password string, err error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", "", errors.New("no Authorization header")
	}
	return parseBasicAuth(auth)
}

func parseBasicAuth(auth string) (username, password string, err error) {
	s1 := strings.SplitN(auth, " ", 2)
	if len(s1) != 2 {
		return "", "", errors.New("failed to parse authentication string")
	}
	if s1[0] != "Basic" {
		return "", "", fmt.Errorf("authorization scheme is %v, not Basic", s1[0])
	}
	c, err := base64.StdEncoding.DecodeString(s1[1])
	if err != nil {
		return "", "", errors.New("failed to parse base64 basic credentials")
	}
	s2 := strings.SplitN(string(c), ":", 2)
	if len(s2) != 2 {
		return "", "", errors.New("failed to parse basic credentials")
	}
	return s2[0], s2[1], nil
}

func (rt *Routes) Login(req *http.Request, r render.Render) {
	email, _, err := basicAuth(req)
	if err != nil {
		r.JSON(500, nil)
	}
	user, err := rt.Db.GetUserByEmail(email)
	if user != nil {
		r.JSON(200, map[string]interface{}{
			"id": user.Id,
		})
	} else {
		r.JSON(500, nil)
	}
}

func (rt *Routes) GetUser(params martini.Params, r render.Render) {
	user, err := rt.Db.GetUser(params["id"])
	if err == nil {
		nConn, _ := rt.Db.GetNumConnections(user.Id)
		r.JSON(200, map[string]interface{}{
			"id":          user.Id,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"description": user.Description,
			"num_degree1": nConn,
			"num_degree2": user.NumDegree2,
		})
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "User not found " + err.Error(),
		})
	}
}

func (rt *Routes) PostUserPicture(params martini.Params, r render.Render, req *http.Request) {
	user, err := rt.Db.GetUser(params["id"])
	if err != nil {
		r.JSON(404, map[string]interface{}{
			"message": "User not found " + err.Error(),
		})
		return
	}
	buf, err := ioutil.ReadAll(req.Body)

	if err != nil {
		r.JSON(404, map[string]interface{}{
			"message": "failed to read picture " + err.Error(),
		})
		return
	}
	err = rt.Db.PostUserPicture(user.Id, buf)
	if err == nil {
		r.JSON(201, map[string]interface{}{"id": user.Id})
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "failed to save picture " + err.Error(),
		})
	}
}

func (rt *Routes) GetUserPicture(params martini.Params, r render.Render, w http.ResponseWriter) {
	user, err := rt.Db.GetUser(params["id"])
	if err != nil {
		r.JSON(404, map[string]interface{}{
			"message": "User not found " + err.Error(),
		})
		return
	}
	image, err := rt.Db.GetUserPicture(user.Id)
	if err == nil {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(image)))
		w.Write(image)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "failed to retrieve picture " + err.Error(),
		})
	}
}

func (rt *Routes) PostUser(r render.Render, req *http.Request) {
	user, err := rt.Db.NewUser(
		"email",
		req.FormValue("email"),
		req.FormValue("password"),
		req.FormValue("first_name"),
		req.FormValue("last_name"),
	)

	if err == nil {
		err = rt.Db.PostUser(user)
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
}
