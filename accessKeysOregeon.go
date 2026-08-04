package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url1         = "https://api.idrivee2.com/api/access_key/add"
	bearerToken1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZDNkODM0ZTQtMzE3MC0xMWYxLTkxNzQtN2NjMjU1ZTUxM2UyIiwiYWRtaW4iOnRydWUsInB1c2VyX2lkIjpudWxsLCJyZXNlbGxlciI6ZmFsc2UsIm1mYV9uYW1lIjpudWxsLCJtZmFfdHlwZSI6bnVsbCwiYWRtaW5faWQiOm51bGwsImlzX3N1YnVzZXJfYWRtaW4iOm51bGwsImhvc3QiOiJjb25zb2xlLmlkcml2ZWUyLmNvbSIsImlhdCI6MTc3NTY0MDMzMCwiZXhwIjoxNzc1NjQwNjMwLCJhdWQiOiJlMi5hcHMuaWRyaXZlIiwiaXNzIjoiZTIuYXBzLmlkcml2ZSIsImp0aSI6ImM4NWEzNjRmLTNkYTMtNWMxYi04ODMyLTFjODlkNTlhMjdiZiJ9.XqVfjQEYczL7u7vKhCc8LmcN5Cq1EDYjKoZOf4Zymvk" // Replace with your actual token
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

// {"name":"sdds","rdns":"h3v1.or8.idrivee2-73.com","permissions":2,"disable_delete_object":false,"disable_delete_version":false,"disable_delete_bucket":false,"expiry_onL":null}
// func sendReq(wg *sync.WaitGroup, id int) {
func sendReq(id int) {
	//defer wg.Done()

	req, err := http.NewRequest("PUT", url1, bytes.NewBuffer(jsonPayLoad1))
	if err != nil {
		fmt.Printf("[Request %d] Error creating the request: %v\n", id, err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken1)
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
	//wg := &sync.WaitGroup{} // ✅ Initialize properly

	//for i := 0; i <= 1; i++ {
	i := 0
	//wg.Add(1)
	//go sendReq(wg, i)
	go sendReq(i)
	time.Sleep(1 * time.Second) // Optional delay
	//}

	//wg.Wait()
	fmt.Println("All requests completed")
}
