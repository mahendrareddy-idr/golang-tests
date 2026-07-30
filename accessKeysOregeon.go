package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	url1           = "https://api.idrivee2.com/api/access_key/add"
	bearerTokenAlt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYmI0OTczZTUtODUzNy0xMWVmLWIxNTgtM2NlY2VmN2NjYTE0IiwiYWRtaW4iOnRydWUsInB1c2VyX2lkIjpudWxsLCJyZXNlbGxlciI6ZmFsc2UsImlhdCI6MTc0ODkzNDUzNiwiZXhwIjoxNzQ4OTM0ODM2LCJhdWQiOiJlMi5hcHMuaWRyaXZlIiwiaXNzIjoiZTIuYXBzLmlkcml2ZSIsImp0aSI6ImE2ZTg5NWFlLTY4OTAtNTQwNy05YTU4LTQ1ZWM5OGFhZDM5MyJ9.T4e9Js4euAYu6x69T-1-Ej0kxLZTkb_3EzmZzRcTtmc"
)

var jsonPayLoad1 = []byte(`{
	"name":"te acce",
	"rdns":"h3v1.or8.idrivee2-73.com",
	"permissions":2,
	"disable_delete_object":false,
	"disable_delete_version":false,
	"disable_delete_bucket":false,
	"expiry_on":null
}`)

func sendReq(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	req, err := http.NewRequest("PUT", url1, bytes.NewBuffer(jsonPayLoad1))
	if err != nil {
		fmt.Printf("[Request %d] Error creating the request: %v\n", id, err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+bearerTokenAlt)
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
	count := 0
	for count < 5 {
		var wg sync.WaitGroup

		for i := 0; i < 2; i++ {
			//wg.Add(1)
			go sendReq(&wg, count*2+i)
		}

		wg.Wait()
		fmt.Printf("Batch %d completed\n", count+1)
		count++
	}
}
