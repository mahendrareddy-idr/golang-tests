package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	apiurl       = "https://api.idrivee2.com/api/service/add"
	bearerToken1 = "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmJlNDU1ZmYtYjk5YS0xMWYwLTkxNzQtN2NjMjU1ZTUxM2UyIiwiYWRtaW4iOnRydWUsInB1c2VyX2lkIjpudWxsLCJyZXNlbGxlciI6ZmFsc2UsIm1mYV9kZXZpY2VfbmFtZSI6bnVsbCwibWZhX3R5cGUiOm51bGwsImlhdCI6MTc2MjMyMDU5OSwiZXhwIjoxNzYyMzIwODk5LCJhdWQiOiJlMi5hcHMuaWRyaXZlIiwiaXNzIjoiZTIuYXBzLmlkcml2ZSIsImp0aSI6IjQ2ODdmOWI5LTNiOGItNTBiZS05ZjRlLTBmNzZjZWM0NzdmYyJ9.ml_Ywkxv0Cyad3_fZdowWcVgN9SSdlFXTPCsDiCPllg"
)

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	client := &http.Client{Timeout: 30 * time.Second}

	for i := 0; i <= 100; i++ {
		//preparing request body
		body := map[string]interface{}{
			"name":                   fmt.Sprintf("100-%d", i),
			"subuser_id":             "6be455ff-b99a-11f0-9174-7cc255e513e2",
			"rdns":                   "u4c4.or.idrivee2-50.com",
			"permissions":            2,
			"disable_delete_object":  false,
			"disable_delete_version": false,
			"disable_delete_bucket":  false,
		}
		jsonBody, _ := json.Marshal(body)

		req, err := http.NewRequest(http.MethodPut, apiurl, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", bearerToken1)
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Request %d failed: %v\n", i, err)
			continue
		}

		respData, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var result struct {
			RespCode int `json:"resp_code"`
			Data     struct {
				User string `json:"user"`
				Pass string `json:"pass"`
			} `json:"data"`
		}

		//var result map[string]interface{}

		if err := json.Unmarshal(respData, &result); err != nil {
			fmt.Printf("Invalid JSON for key %d: %v\nResponse: %s\n", i, err, string(respData))
			continue
		}

		/*user, userOk := result["user"].(string)
		pass, passOk := result["pass"].(string)
		*/

		//if !userOk || !passOk {
		if result.RespCode != 0 || result.Data.User == "" || result.Data.Pass == "" {
			fmt.Println("Missing fields in response for key %d: %s\n", i, string(respData))
			time.Sleep(2 * time.Second)
			continue
		}

		//Write to file line
		line := fmt.Sprintf("generated key %d: %s\n", result.Data.User, result.Data.Pass)

		file.WriteString(line)
		fmt.Println("Generated key %d: %s\n", i, line)
	}
	fmt.Println("100 keys generated and saved to output_keys.txt")
}
