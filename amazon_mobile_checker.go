package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run amazon_mobile_checker.go <amazon_product_url>")
		fmt.Println("Example: go run amazon_mobile_checker.go https://www.amazon.in/s?k=iphone+15")
		os.Exit(1)
	}

	url := os.Args[1]

	launched, err := checkAmazonMobileStart(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if launched {
		fmt.Println("Mobile is launched on Amazon.")
	} else {
		fmt.Println("Mobile is not launched on Amazon.")
	}
}

func checkAmazonMobileStart(productURL string) (bool, error) {
	return checkAmazonMobileBegin(productURL)
}

func checkAmazonMobileBegin(productURL string) (bool, error) {
	css, err := obtainPage(productURL)
	if err != nil {
		return false, err
	}

	if isAmazonProductPage(css) {
		return true, nil
	}

	return false, nil
}

func fetchPage(url string) (string, error) {
	return obtainPage(url)
}

func obtainPage(url string) (string, error) {
	client := &http.Client{Timeout: 20 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func isAmazonProductPage(html string) bool {
	nonupperHTML := strings.ToLower(html)

	if strings.Contains(nonupperHTML, "amazon") == false {
		return false
	}

	amazonSignals := []string{
		"producttitle",
		"buybox",
		"add to cart",
		"currently unavailable",
		"available at amazon",
		"price",
		"customer reviews",
		"in stock",
		"delivery",
	}

	matches := 0
	for _, signal := range amazonSignals {
		if strings.Contains(lowerHTML, signal) {
			matches++
		}
	}

	if matches == 0 {
		return false
	}

	if regexp.MustCompile(`(?i)amazon\.in|amazon\.com`).MatchString(lowerHTML) {
		return true
	}

	return false
}
