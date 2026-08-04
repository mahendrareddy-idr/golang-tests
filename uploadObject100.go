package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

const (
	inputFile   = "100accesskeys.txt"
	resultsFile = "results1.txt"
	endpointUrl = "https://s3.us-west-1.idrivee2.com"
	objectKey   = `verify.csv`
	//objectKey1  = `verify2.csv`
	bucketName = "100times-test"
)

var (
	reUser = regexp.MustCompile(`user"\s*:\s*"([^"]+)"`)
	rePass = regexp.MustCompile(`pass"\s*:\s*"([^"]+)"`)
)

func main() {
	// open input file
	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", inputFile, err)
		return
	}
	defer f.Close()

	// create a output file
	out, err := os.Create(resultsFile)
	if err != nil {
		fmt.Printf("Failed to create a file %s: %v", resultsFile, err)
		return
	}
	defer out.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		user := extract(reUser, line)
		pass := extract(rePass, line)

		if user == "" || pass == "" {
			msg := fmt.Sprintf("Line %d: could not parse credentials. Skipping.\n", lineNum)
			fmt.Println(msg)
			out.WriteString(msg)
			continue
		}
		wg.Add(1)
		go func(user, pass string, lineNum int) {
			defer wg.Done()
			header := fmt.Sprintf("\n ===============\n  Line %d\nAccessKey: %s\nSecretKey: %s\n=============================\n", lineNum, user, pass)
			fmt.Println(header)
			mu.Lock()
			out.WriteString(header)
			mu.Unlock()

			args := []string{
				"s3api", "put-object",
				"--endpoint-url", endpointUrl,
				"--key", objectKey,
				"--bucket", bucketName,
				"--body", objectKey,
			}
			stdout, stderr, err := runAwsWithEnv(user, pass, args)
			block := fmt.Sprintf("Output for Line %d:\nSTDOUT:\n%s\nSTDERR:\n%s\n", lineNum, stdout, stderr)
			fmt.Print(block)
			mu.Lock()
			out.WriteString(block)
			mu.Unlock()

			status := "Success"
			if err != nil {
				status = fmt.Sprintf("failed: %v", err)
			}
			resultLine := fmt.Sprintf("Completed Line %d — Status: %s\n", lineNum, status)
			fmt.Print(resultLine)
			mu.Lock()
			out.WriteString(resultLine)
			mu.Unlock()
		}(user, pass, lineNum)
	}
	wg.Wait()
	fmt.Printf("\n All done! Full results written to %s\n", resultsFile)
}

// extract helper

func extract(r *regexp.Regexp, s string) string {
	m := r.FindStringSubmatch(s)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func runAwsWithEnv(accessKey, secretKey string, args []string) (string, string, error) {
	cmd := exec.Command("aws", args...)
	env := os.Environ()
	env = append(env,
		"AWS_ACCESS_KEY_ID="+accessKey,
		"AWS_SECRET_ACCESS_KEY="+secretKey,
		"AWS_DEFAULT_REGION=us-west-1",
	)
	cmd.Env = env
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	return strings.TrimSpace(outb.String()), strings.TrimSpace(errb.String()), err
}
