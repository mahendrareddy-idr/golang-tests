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

const (
	defaultSourceURL  = "https://www.google.com/search?q=telugu+movies+new+releases+this+week"
	defaultOutputFile = "telugu_movies_this_week.txt"
)

func main() {
	sourceURL := defaultSourceURL
	if len(os.Args) > 1 {
		sourceURL = os.Args[1]
	}

	outputFile := defaultOutputFile
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	movies, err := getBranchAMovies(sourceURL)
	if err != nil {
		fmt.Println("Failed to fetch release data:", err)
		os.Exit(1)
	}

	if len(movies) == 0 {
		fmt.Println("No Telugu movie release titles were found from the source.")
		return
	}

	if err := saveMoviesToFile(outputFile, movies); err != nil {
		fmt.Println("Failed to write file:", err)
		os.Exit(1)
	}

	fmt.Println("Saved Telugu movie releases to:", outputFile)
	for _, movie := range movies {
		fmt.Println("- " + movie)
	}
}

func getBranchAMovies(sourceURL string) ([]string, error) {
	html, err := fetchWebPage(sourceURL)
	if err != nil {
		return nil, err
	}

	movies := extractMovieTitles(html)
	if len(movies) == 0 {
		return nil, fmt.Errorf("no movie titles could be extracted from the web page")
	}

	return movies, nil
}

func fetchWebPage(url string) (string, error) {
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func extractMovieTitles(html string) []string {
	titlePattern := regexp.MustCompile(`(?is)<h[1-3][^>]*>(.*?)</h[1-3]>|(?is)<a[^>]+href=[^>]*>(.*?)</a>`)
	matches := titlePattern.FindAllStringSubmatch(html, 200)

	seen := map[string]bool{}
	results := []string{}

	for _, match := range matches {
		for _, part := range match[1:] {
			clean := cleanText(part)
			if clean == "" || strings.Contains(strings.ToLower(clean), "google") || strings.Contains(strings.ToLower(clean), "search") {
				continue
			}
			if looksLikeMovieTitle(clean) && !seen[clean] {
				seen[clean] = true
				results = append(results, clean)
			}
		}
	}

	return results
}

func cleanText(raw string) string {
	text := strings.ReplaceAll(raw, "&amp;", "&")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "<b>", "")
	text = strings.ReplaceAll(text, "</b>", "")

	re := regexp.MustCompile(`(?is)<[^>]+>`)
	text = re.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func looksLikeMovieTitle(value string) bool {
	if len(value) < 3 || len(value) > 120 {
		return false
	}

	lower := strings.ToLower(value)
	blocked := []string{
		"signin", "login", "privacy", "terms", "settings", "images", "videos",
		"news", "home", "all", "about", "contact", "cookies", "advertise",
		"results", "next", "previous", "google", "search",
	}
	for _, word := range blocked {
		if strings.Contains(lower, word) {
			return false
		}
	}

	return regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 .:'’&()-]+$`).MatchString(value)
}

func saveMoviesToFile(path string, movies []string) error {
	content := strings.Builder{}
	for i, movie := range movies {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, movie))
	}

	return os.WriteFile(path, []byte(content.String()), 0644)
}
