package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

const (
	totalBuckets = 40000
	workers      = 20 // Tune this (10–50 is usually safe)
)

func configureAWS() {
	os.Setenv("AWS_ACCESS_KEY_ID", "UqMoxHJ8vNc4vrvTp5ch")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "hb8JEjkzmixK9RgFE4nNEde2mhKnwgaRkBbsGlnV")
	os.Setenv("AWS_DEFAULT_REGION", "us-west-1")
}

func createBucket(bucketName string, endpointURL string) error {
	cmd := exec.Command(
		"aws", "s3api", "create-bucket",
		"--bucket", bucketName,
		"--region", "us-east-1",
		"--create-bucket-configuration", "LocationConstraint=us-east-1",
		"--endpoint-url", endpointURL,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", bucketName, output)
	}

	fmt.Printf("✅ %s created\n", bucketName)
	return nil
}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup, endpoint string) {
	defer wg.Done()

	for bucket := range jobs {
		if err := createBucket(bucket, endpoint); err != nil {
			log.Printf("❌ Worker %d error: %v\n", id, err)
		}
	}
}

func main() {
	configureAWS()

	endpoint := "https://s3.us-west-1.idrivee2.com"

	jobs := make(chan string, totalBuckets)
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg, endpoint)
	}

	// Send jobs
	for i := 1; i <= totalBuckets; i++ {
		bucketName := fmt.Sprintf("mahendra2-%d", i)
		jobs <- bucketName
	}
	close(jobs)

	// Wait for completion
	wg.Wait()

	fmt.Println("🎉 All buckets processed")
}
