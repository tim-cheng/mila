package main

import (
	"bytes"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func testClientJsonResp(email, password, method, path, params string) (retStatus int, retResp *simplejson.Json, err error) {
	retStatus = 0
	retResp = nil
	err = nil

	client := &http.Client{}
	req, err := http.NewRequest(method, "http://localhost:8080"+path, bytes.NewBufferString(params))
	if err != nil {
		return
	}
	if email != "" {
		req.SetBasicAuth(email, password)
	}
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	retStatus = resp.StatusCode
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// ignore json decode error
	retResp, _ = simplejson.NewJson(buf)
	return
}

func testClient(email, password, method, path, params string) int {
	code, _, err := testClientJsonResp(email, password, method, path, params)
	if err != nil {
		code = 0
	}
	return code
}

func testPostImage(email, password, path, filename string) int {
	client := &http.Client{}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0
	}

	req, err := http.NewRequest("POST", "http://localhost:8080"+path, bytes.NewReader(buf))
	if err != nil {
		return 0
	}
	if email != "" {
		req.SetBasicAuth(email, password)
	}
	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Content-Length", strconv.Itoa(len(buf)))
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	return resp.StatusCode
}

func createUser(email, password, firstname, lastname string) int {
	params := url.Values{}
	params.Add("email", email)
	params.Add("password", password)
	params.Add("first_name", firstname)
	params.Add("last_name", lastname)
	return testClient("", "", "POST", "/users", params.Encode())
}

func createKid(pid, name, birthday, boygirl string) int {
	params := url.Values{}
	params.Add("name", name)
	params.Add("birthday", birthday)
	params.Add("type", boygirl)
	return testClient("", "", "POST", "/users/"+pid+"/kids", params.Encode())
}

func createConnection(email, password, id1, id2 string) int {
	params := url.Values{}
	params.Add("user1_id", id1)
	params.Add("user2_id", id2)
	return testClient(email, password, "POST", "/connections", params.Encode())
}

func createPost(email, password, userId, content string) int {
	params := url.Values{}
	params.Add("user_id", userId)
	params.Add("body", content)
	return testClient(email, password, "POST", "/posts", params.Encode())
}

func createComment(email, password, userId, postId, content string) int {
	params := url.Values{}
	params.Add("user_id", userId)
	params.Add("body", content)
	return testClient(email, password, "POST", "/posts/"+postId+"/comments", params.Encode())
}

func testStar(method, email, password, userId, postId string) int {
	params := url.Values{}
	params.Add("user_id", userId)
	return testClient(email, password, method, "/posts/"+postId+"/stars?"+params.Encode(), "")
}

func createStar(email, password, userId, postId string) int {
	return testStar("PUT", email, password, userId, postId)
}

func deleteStar(email, password, userId, postId string) int {
	return testStar("DELETE", email, password, userId, postId)
}

func checkCode(t *testing.T, msg string, code int, expect int) {
	t.Log("test: " + msg)
	if code == expect {
		t.Log("passed!")
	} else {
		t.Errorf("failed: expect %d, got %d\n", expect, code)
	}
}

func checkCode2(t *testing.T, msg string, code int, code2 int, expect int, expect2 int) {
	t.Log("test: " + msg)
	if code == expect && code2 == expect2 {
		t.Log("passed!")
	} else {
		t.Errorf("failed: expect (%d, %d), got (%d, %d)\n", expect, expect2, code, code2)
	}
}

func TestBasic(t *testing.T) {
	go startServer()
	time.Sleep(time.Second * 2)

	// create users
	numList := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, v := range numList {
		e, p, f, l := "user"+v+"@test.com", "testtest"+v, "First"+v, "Last"+v
		checkCode(t, "create user", createUser(e, p, f, l), 201)
	}

	checkCode(t, "create user with same username", createUser("user4@test.com", "sdfsfss", "sfsfs", "sfsfs"), 404)

	e, p := "user1@test.com", "testtest1"

	checkCode(t, "get user info self", testClient(e, p, "GET", "/users/1", ""), 200)
	checkCode(t, "get user info others", testClient(e, p, "GET", "/users/2", ""), 200)
	checkCode(t, "get user info doesn't exist", testClient(e, p, "GET", "/users/100", ""), 404)
	checkCode(t, "get user info auth fail", testClient(e, "testtest", "GET", "/users/1", ""), 401)

	checkCode(t, "create connection", createConnection(e, p, "1", "2"), 201)
	checkCode(t, "create connection", createConnection(e, p, "1", "3"), 201)
	checkCode(t, "create connection", createConnection(e, p, "1", "4"), 201)
	checkCode(t, "create connection", createConnection(e, p, "5", "1"), 201)
	checkCode(t, "create connection", createConnection(e, p, "6", "1"), 201)

	checkCode(t, "can't connect again", createConnection(e, p, "1", "2"), 404)
	checkCode(t, "can't connect with self", createConnection(e, p, "1", "1"), 404)
	checkCode(t, "connect auth fail", createConnection(e, "testtest", "1", "7"), 401)

	checkCode(t, "create post", createPost(e, p, "1", "This is post1"), 201)
	checkCode(t, "create post", createPost(e, p, "1", "This is post2"), 201)
	checkCode(t, "create post", createPost(e, p, "1", "This is post3"), 201)
	checkCode(t, "create post", createPost(e, p, "1", "This is post4"), 201)
	checkCode(t, "create post", createPost(e, p, "2", "This is post5"), 201)
	checkCode(t, "create post", createPost(e, p, "2", "This is post6"), 201)
	checkCode(t, "delete post", testClient(e, p, "DELETE", "/posts/6", ""), 200)
	// TODO: test delete post permission
	checkCode(t, "create post auth", createPost(e, "testtest", "1", "This is post2"), 401)
	checkCode(t, "create post non-existent user", createPost(e, p, "100", "This is post2"), 404)

	checkCode(t, "create comment", createComment(e, p, "1", "1", "This is comment1"), 201)
	checkCode(t, "create comment", createComment(e, p, "1", "2", "This is comment2"), 201)
	checkCode(t, "create comment", createComment(e, p, "2", "1", "This is comment3"), 201)
	checkCode(t, "create comment", createComment(e, p, "2", "2", "This is comment4"), 201)
	checkCode(t, "create comment non-existent user", createComment(e, p, "100", "1", "This is comment1"), 404)
	checkCode(t, "create comment non-existent post", createComment(e, p, "1", "100", "This is comment1"), 404)
	checkCode(t, "create comment auth", createComment(e, "testtest", "1", "1", "This is comment1"), 401)

	checkCode(t, "get comment, has comments", testClient(e, p, "GET", "/posts/1/comments", ""), 200)
	checkCode(t, "get comment, no comments", testClient(e, p, "GET", "/posts/3/comments", ""), 404)
	checkCode(t, "get comment, no posts", testClient(e, p, "GET", "/posts/100/comments", ""), 404)

	checkCode(t, "create star", createStar(e, p, "1", "1"), 200)
	checkCode(t, "create star non-existent user", createStar(e, p, "100", "1"), 404)
	checkCode(t, "create star non-existent post", createStar(e, p, "1", "100"), 404)
	checkCode(t, "create star auth", createStar(e, "testtest", "1", "2"), 401)
	checkCode(t, "can't start again", createStar(e, p, "1", "1"), 404)

	checkCode(t, "delete star", deleteStar(e, p, "1", "1"), 200)
	checkCode(t, "delete non-existent star", deleteStar(e, p, "1", "2"), 404)

	checkCode(t, "create star", createStar(e, p, "1", "1"), 200)
	checkCode(t, "create star", createStar(e, p, "2", "1"), 200)
	checkCode(t, "create star", createStar(e, p, "1", "2"), 200)
	checkCode(t, "create star", createStar(e, p, "2", "2"), 200)

	checkCode(t, "login", testClient(e, p, "GET", "/login", ""), 200)
	checkCode(t, "login auth fail", testClient(e, "testtest", "GET", "/login", ""), 401)

	checkCode(t, "get posts", testClient(e, p, "GET", "/posts?user_id=1", ""), 200)
	checkCode(t, "get posts without user_id", testClient(e, p, "GET", "/posts", ""), 404)
	checkCode(t, "get posts with invalid user_id", testClient(e, p, "GET", "/posts?user_id=abc", ""), 404)

	checkCode(t, "post profile picture", testPostImage(e, p, "/users/1/picture", "./tests/sherry.jpg"), 201)
	checkCode(t, "post profile picture", testPostImage(e, p, "/users/2/picture", "./tests/tim.jpg"), 201)
	checkCode(t, "get profile picture", testClient(e, p, "GET", "/users/1/picture", ""), 200)
	checkCode(t, "get profile picture", testClient(e, p, "GET", "/users/2/picture", ""), 200)

	checkCode(t, "post post picture", testPostImage(e, p, "/posts/1/picture", "./tests/post1.jpg"), 201)
	checkCode(t, "get post picture", testClient(e, p, "GET", "/posts/1/picture", ""), 200)

	// TODO: test FB credential handling
	//checkCode(t, "login with fb", testClient("fb_user@test.com", "random_access_token", "GET", "/login_facebook?first_name=Fb&last_name=User", ""), 201)
	//checkCode(t, "login again with fb", testClient("fb_user@test.com", "random_access_token", "GET", "/login_facebook?first_name=Fb&last_name=User", ""), 200)

	// search
	checkCode(t, "search users", testClient(e, p, "GET", "/users", ""), 404)
	checkCode(t, "search users", testClient(e, p, "GET", "/users?search=blah", ""), 404)
	checkCode(t, "search users", testClient(e, p, "GET", "/users?search=First3", ""), 200)
	checkCode(t, "search users", testClient(e, p, "GET", "/users?search=Last3", ""), 200)
	checkCode(t, "search users", testClient(e, p, "GET", "/users?search=First3%20Last3", ""), 200)

	checkCode(t, "create user", createUser("search1@test.com", "searh1", "First3", "Last3"), 201)
	checkCode(t, "create user", createUser("search2@test.com", "searh1", "First3", "Last3"), 201)
	checkCode(t, "create user", createUser("search3@test.com", "searh1", "First3", "Last3"), 201)
	checkCode(t, "search users", testClient(e, p, "GET", "/users?search=First3", ""), 200)

	checkCode(t, "invite user", testClient(e, p, "POST", "/users/1/invite/7", ""), 201)
	checkCode(t, "invite user", testClient(e, p, "POST", "/users/1/invite/8", ""), 201)
	checkCode(t, "get invites 200", testClient(e, p, "GET", "/connections/1", ""), 200)
	checkCode(t, "get invites 404", testClient(e, p, "GET", "/connections/8", ""), 404)
	checkCode(t, "get connections 200", testClient(e, p, "GET", "/connections/1", ""), 200)
	checkCode(t, "get connections 404", testClient(e, p, "GET", "/connections/8", ""), 404)

	checkCode(t, "accept invite", testClient(e, p, "DELETE", "/users/1/invite/8", ""), 200)
	checkCode(t, "new connection established", testClient(e, p, "GET", "/connections/8", ""), 200)

	checkCode(t, "get activities", testClient(e, p, "GET", "/activities/8", ""), 404)
	checkCode(t, "create post", createPost(e, p, "1", "This is new post"), 201)
	checkCode(t, "create post", createPost(e, p, "8", "This is post8"), 201)
	checkCode(t, "create comment", createComment(e, p, "1", "9", "This is my comment"), 201)
	checkCode(t, "create star", createStar(e, p, "1", "9"), 200)

	// actiivities are not added right away
	time.Sleep(time.Second * 1)

	code, resp, _ := testClientJsonResp(e, p, "GET", "/activities/8", "")
	respAry, _ := resp.Array()
	fmt.Println("response ", resp)
	checkCode2(t, "get activities post", code, len(respAry), 200, 3)

	// kids
	checkCode(t, "add kids", createKid("8", "jane", "2010-03-05", "girl"), 201)
	checkCode(t, "add kids", createKid("8", "toms", "2002-04-05", "boy"), 201)
	checkCode(t, "get kids", testClient(e, p, "GET", "/users/8/kids", ""), 200)
	checkCode(t, "get kids", testClient(e, p, "GET", "/users/7/kids", ""), 404)
	checkCode(t, "get kids", testClient(e, p, "DELETE", "/users/8/kids/1", ""), 200)
	checkCode(t, "get kids", testClient(e, p, "DELETE", "/users/8/kids/2", ""), 200)
	checkCode(t, "get kids", testClient(e, p, "GET", "/users/8/kids", ""), 404)

	checkCode(t, "add kids", createKid("8", "jane", "2010-03-05", "girl"), 201)
	checkCode(t, "add kids", createKid("8", "toms", "2002-04-05", "boy"), 201)

	checkCode(t, "get posts", testClient(e, p, "GET", "/posts?user_id=1&degree=1", ""), 200)
	checkCode(t, "get posts", testClient(e, p, "GET", "/posts?user_id=1&degree=2", ""), 200)
	checkCode(t, "get posts", testClient(e, p, "GET", "/posts?user_id=8&degree=2", ""), 200)
	checkCode(t, "create post", createPost(e, p, "2", "This is another post"), 201)

	// feeds are not added right away
	time.Sleep(time.Second * 1)
	checkCode(t, "get posts", testClient(e, p, "GET", "/posts?user_id=8&degree=2", ""), 200)

}
