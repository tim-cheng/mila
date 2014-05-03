package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func testClient(email, password, method, path, params string) {
	client := &http.Client{}
	req, err := http.NewRequest(method, "http://localhost:8080"+path, bytes.NewBufferString(params))
	req.SetBasicAuth(email, password)
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	}
	resp, err := client.Do(req)
	fmt.Printf("req: %s\nresp: %v\nerr: %v\n", path, resp, err)
	contents, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println(string(contents))
	}
}

func TestBasic(t *testing.T) {
	go startServer()
	time.Sleep(time.Second * 2)

	e, p := "a@b.c", "abcd1234"

	testClient(e, p, "GET", "/users/1", "")

	params := url.Values{}
	params.Add("email", "test@test.com")
	params.Add("password", "234!G)...")
	params.Add("first_name", "fslfs")
	params.Add("last_name", "fasfsf")
	testClient(e, p, "POST", "/users", params.Encode())

}
