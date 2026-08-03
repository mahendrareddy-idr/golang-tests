package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//const (
	url         = "https://api.idrivee2.com/api/access_key/add"
	bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZDNkODM0ZTQtMzE3MC0xMWYxLTkxNzQtN2NjMjU1ZTUxM2UyIiwiYWRtaW4iOnRydWUsInB1c2VyX2lkIjpudWxsLCJyZXNlbGxlciI6ZmFsc2UsIm1mYV9uYW1lIjpudWxsLCJtZmFfdHlwZSI6bnVsbCwiYWRtaW5faWQiOm51bGwsImlzX3N1YnVzZXJfYWRtaW4iOm51bGwsImhvc3QiOiJjb25zb2xlLmlkcml2ZWUyLmNvbSIsImlhdCI6MTc3NTY0MDcxNSwiZXhwIjoxNzc1NjQxMDE1LCJhdWQiOiJlMi5hcHMuaWRyaXZlIiwiaXNzIjoiZTIuYXBzLmlkcml2ZSIsImp0aSI6ImM4NWEzNjRmLTNkYTMtNWMxYi04ODMyLTFjODlkNTlhMjdiZiJ9.GLM0ZbdQp87uPFNd2iVhMnVzNvNWh84u4XkLAJgXArw"
//)

var jsonPayLoad = []byte(`{
"name":"123",
"rdns":"h3v1.or8.idrivee2-73.com",
"permissions":2,
"disable_delete_object":false,
"disable_delete_version":false,
"disable_delete_bucket":false,
"expiry_onL":null}`)

func sendRequest(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayLoad))
	if err != nil {
		fmt.Printf("[Request %d] Error creating the request: %v\n", id, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[Request %d] Request failed: %v\n", id, err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("[Request %d] Status: %s, Response: %s\n", id, resp.Status, string(body))
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i <= 1; i++ {
		wg.Add(1)
		time.Sleep(60)
		go sendRequest(&wg, i)
	}
	wg.Wait()
	fmt.Println("All 250 requests completed")
}
