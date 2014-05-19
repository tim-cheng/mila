package routes

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/mostafah/mandrill"
	"io"
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

func (rt *Routes) downloadFacebookPicture(userId int64, fbId string) {
	// upload facebook picture
	res, err := http.Get("https://graph.facebook.com/" + fbId + "/picture?type=small")
	if err == nil {
		buf, err := ioutil.ReadAll(res.Body)
		if err == nil {
			err = rt.Db.PostUserPicture(userId, buf)
			fmt.Printf("upload fb picture, buflen= %d, err=%v\n", len(buf), err)
		}
	}
}

func getFbIdFromToken(token string) (string, error) {
	res, err := http.Get("https://graph.facebook.com/me?fields=id&access_token=" + token)
	if err != nil {
		return "", err
	}
	var msg struct {
		Id string `json:"id"`
	}
	d := json.NewDecoder(res.Body)
	err = d.Decode(&msg)
	if err != nil && err != io.EOF {
		return "", err
	}
	if msg.Id == "" {
		return "", errors.New("FB auth failed")
	}
	return msg.Id, nil
}

func (rt *Routes) LoginFacebook(r render.Render, req *http.Request) {

	retStatus := 401
	// break out for failure if anything failed
	for {
		email, password, err := basicAuth(req)
		if err != nil {
			break
		}
		fmt.Printf("auth fb: %v, %v\n", email, password)

		fbId, err := getFbIdFromToken(password)
		if err != nil {
			break
		}

		user, err := rt.Db.GetUserByEmail(email + "@fb")
		if user != nil && err == nil {
			// user exist
			retStatus = 200
			if fbId != req.FormValue("fb_id") {
				break
			}
			err = rt.Db.UpdatePassword(user, password)
			if err != nil {
				break
			}
		} else {
			// user doesn't exist
			retStatus = 201
			user, err = rt.Db.NewUser(
				"facebook",
				email+"@fb",
				password,
				req.FormValue("first_name"),
				req.FormValue("last_name"),
				req.FormValue("fb_id"),
			)
			if err != nil {
				break
			}
			err = rt.Db.PostUser(user)
			if err != nil {
				break
			}
			go rt.downloadFacebookPicture(user.Id, req.FormValue("fb_id"))
			go sendNewUserEmail(email, user.FirstName)
		}

		r.JSON(retStatus, map[string]interface{}{
			"id": user.Id,
		})
		return
	}
	r.JSON(retStatus, map[string]interface{}{
		"message": "Not authorized",
	})
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
		"",
	)

	if err == nil {
		err = rt.Db.PostUser(user)
	}
	if err == nil {
		go sendNewUserEmail(user.Email, user.FirstName)
		r.JSON(201, map[string]interface{}{
			"id": user.Id,
		})

	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to add user " + err.Error(),
		})
	}
}

func (rt *Routes) SearchUsers(r render.Render, req *http.Request) {
	search := req.FormValue("search")
	if search == "" {
		r.JSON(404, map[string]interface{}{
			"message": "search param missing",
		})
		return
	}

	words := strings.Split(search, " ")
	if len(words) > 2 {
		r.JSON(404, map[string]interface{}{
			"message": "more than 2 words in query",
		})
		return
	}

	var res []int64
	if len(words) == 1 {
		// contact results from first name / last name search
		fRes, _ := rt.Db.GetUserIdByFirstName(words[0])
		lRes, _ := rt.Db.GetUserIdByLastName(words[0])
		res = append(fRes, lRes...)
	} else {
		// first name last name
		res, _ = rt.Db.GetUserIdByFullName(words[0], words[1])
	}

	if len(res) == 0 {
		r.JSON(404, map[string]interface{}{
			"message": "no users found",
		})
		return
	}
	resMsg := make([]map[string]int64, len(res), len(res))
	for i, id := range res {
		resMsg[i] = map[string]int64{"id": id}
	}

	r.JSON(200, resMsg)
}

func sendNewUserEmail(email, first_name string) {
	mandrill.Key = "izQqlSTrNP4ZKZQ_rtM3-Q"
	msg := mandrill.NewMessageTo(email, first_name)
	msg.HTML = "<p>Welcome to Parent2D</p>"
	msg.Text = "Welcome to Parent2D" // optional
	msg.Subject = "Welcome to Parent2D"
	msg.FromEmail = "sherry@parent2d.com"
	msg.FromName = "Parent2D"
	res, err := msg.Send(false)
	fmt.Printf("res = %v, err = %v\n", res, err)
}
