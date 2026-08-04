package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseItem represents the structure of each item in the response JSON array
type ResponseItem struct {
	KeyId string `json:"KeyId"`
	Rdns  string `json:"rdns"`
}

func main() {
	// API endpoints
	initialURL := "https://api.idrivee2.com/api/"
	removeURL := "https://api.idrivee2.com/api/service/remove"

	// Bearer token
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYmI0OTczZTUtODUzNy0xMWVmLWIxNTgtM2NlY2VmN2NjYTE0IiwiYWRtaW4iOnRydWUsInB1c2VyX2lkIjpudWxsLCJyZXNlbGxlciI6ZmFsc2UsImlhdCI6MTc0ODkzNDUzNiwiZXhwIjoxNzQ4OTM0ODM2LCJhdWQiOiJlMi5hcHMuaWRyaXZlIiwiaXNzIjoiZTIuYXBzLmlkcml2ZSIsImp0aSI6ImE2ZTg5NWFlLTY4OTAtNTQwNy05YTU4LTQ1ZWM5OGFhZDM5MyJ9.T4e9Js4euAYu6x69T-1-Ej0kxLZTkb_3EzmZzRcTtmc"

	// Step 1: Send initial POST request with empty JSON payload
	jsonData := []byte(`{}`)

	req, err := http.NewRequest("POST", initialURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating initial request:", err)
		return
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending initial request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var items []ResponseItem
	err = json.Unmarshal(body, &items)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Println(string(body))
		return
	}

	// Step 2: Loop through each item and send POST to /service/remove
	for _, item := range items {
		removePayload := map[string]string{
			"svc":  item.KeyId,
			"rdns": item.Rdns,
		}
		payloadBytes, _ := json.Marshal(removePayload)

		removeReq, err := http.NewRequest("POST", removeURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error creating remove request:", err)
			continue
		}
		removeReq.Header.Set("Authorization", "Bearer "+bearerToken)
		removeReq.Header.Set("Content-Type", "application/json")

		removeResp, err := client.Do(removeReq)
		if err != nil {
			fmt.Println("Error sending remove request for KeyId", item.KeyId, ":", err)
			continue
		}
		defer removeResp.Body.Close()

		// Optionally print status or response
		fmt.Printf("Removed KeyId: %s, Response: %s\n", item.KeyId, removeResp.Status)
	}
}
