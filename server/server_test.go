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

	//	testClient(e, p, "GET", "/users/1", "")

	var code int
	// create users
	numList := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, v := range numList {
		e, p, f, l := "user"+v+"@test.com", "testtest"+v, "First"+v, "Last"+v
		code = createUser(e, p, f, l)
		checkCode(t, "create user", code, 201)
	}

	// create user with same username:
	code = createUser("user4@test.com", "sdfsfss", "sfsfs", "sfsfs")
	checkCode(t, "create user with same username", code, 404)

	// get user
	e, p := "user1@test.com", "testtest1"
	code = testClient(e, p, "GET", "/users/1", "")
	checkCode(t, "get user info self", code, 200)
	code = testClient(e, p, "GET", "/users/2", "")
	checkCode(t, "get user info others", code, 200)

	// user doesn't exist
	code = testClient(e, p, "GET", "/user/100", "")
	checkCode(t, "get user info doesn't exist", code, 404)

	// auth fail
	testClient(e, "testtest", "GET", "/user/1", "")
	checkCode(t, "get user info auth fail", code, 404)
}
