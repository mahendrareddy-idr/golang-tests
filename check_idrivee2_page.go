package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := "https://console.idrivee2.com"

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Response looks successful.")
	} else {
		fmt.Println("Response returned an error status.")
	}
}
