package main

import (
  "testing"
  "net/http"
  "fmt"
  "time"
)

func testGet(path string) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", "http://localhost:8080"+path, nil)
  req.SetBasicAuth("a@b.c", "abcd1234")
  resp, err := client.Do(req)
  fmt.Printf("req: %s\nresp: %s\nerr: %v\n", path, resp, err)
}

func TestBasic(t *testing.T) {
  go startServer()
  time.Sleep(time.Second * 2)

  testGet("/users/1")
}
