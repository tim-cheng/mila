package main

import (
	"bytes"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func testClient(email, password, method, path, params string) int {
	client := &http.Client{}
	req, err := http.NewRequest(method, "http://localhost:8080"+path, bytes.NewBufferString(params))
	if err != nil {
		return 0
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
		return 0
	}
	return resp.StatusCode
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

	checkCode(t, "login with fb", testClient("fb_user@test.com", "random_access_token", "GET", "/login_facebook?first_name=Fb&last_name=User", ""), 201)
	checkCode(t, "login again with fb", testClient("fb_user@test.com", "random_access_token", "GET", "/login_facebook?first_name=Fb&last_name=User", ""), 200)
}
