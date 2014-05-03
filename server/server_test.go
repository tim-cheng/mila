package main

import (
	"bytes"
	//"fmt"
	//"io/ioutil"
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
	//fmt.Printf("req: %s\nresp: %v\nerr: %v\n", path, resp, err)
	//contents, err := ioutil.ReadAll(resp.Body)
	//if err == nil {
	//  fmt.Println(string(contents))
	//}

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

}
