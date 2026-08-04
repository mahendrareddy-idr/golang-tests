package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	endpoint    = "https://b6w3.or4.idrivee2-65.com"
	bucket      = "abr-test"
	key         = "internal-battery-inside-deck.png"
	numRequests = 1900000
	concurrency = 200 // Number of concurrent workers
)

func main() {
	start := time.Now()

	// Load AWS config with custom endpoint
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: "us-east-1",
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			}),
		),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	var wg sync.WaitGroup
	requests := make(chan int, concurrency)

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for range requests {
				check, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
				})
				if err != nil {
					var nfe *types.NotFound
					if errors.As(err, &nfe) {
						log.Printf("[Worker %d] Object not found", workerID)
					} else {
						log.Printf("[Worker %d] Error: %v", workerID, err)
					}
				} else {
					output, err := json.MarshalIndent(check, "", "  ")
					if err != nil {
						log.Printf("[Worker %d] Failed to marshal output: %v", workerID, err)
					} else {
						fmt.Printf("[Worker %d] HeadObject result:\n%s\n", workerID, string(output))
					}
				}
			}
		}(i)
	}

	// Send jobs to workers
	for i := 0; i < numRequests; i++ {
		requests <- i
		fmt.Println("Current iteration is ", i)
	}
	close(requests)

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Completed %d head-object requests in %s\n", numRequests, elapsed)
}
