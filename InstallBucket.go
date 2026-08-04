package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func configureAWS() error {
	// Setting AWS environment variables directly (equivalent to 'aws configure')
	// Replace these with your actual AWS credentials and region
	os.Setenv("AWS_ACCESS_KEY_ID", "UqMoxHJ8vNc4vrvTp5ch")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "hb8JEjkzmixK9RgFE4nNEde2mhKnwgaRkBbsGlnV")
	os.Setenv("AWS_DEFAULT_REGION", "us-west-1")

	// Confirm that the environment variables are set
	fmt.Println("AWS Configuration set with access key, secret key, and region.")
	return nil
}

func createBucket(bucketName string, endpointURL string) error {
	// Include the --endpoint-url argument
	cmd := exec.Command("aws", "s3api", "create-bucket", "--bucket", bucketName,
		"--region", "us-east-1",
		"--create-bucket-configuration", "LocationConstraint=us-east-1",
		"--endpoint-url", endpointURL)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %s", bucketName, output)
	}
	fmt.Printf("Bucket %s created successfully.\n", bucketName)
	return nil
}

func main() {
	// Step 1: Configure AWS CLI with environment variables (equivalent to aws configure)
	if err := configureAWS(); err != nil {
		log.Fatalf("AWS configuration failed: %v\n", err)
	}

	// Set your custom endpoint URL here
	endpointURL := "https://s3.us-west-1.idrivee2.com" // Replace with your actual endpoint URL

	// Step 2: Loop to create 1000 buckets
	for i := 1; i <= 10000; i++ {
		bucketName := fmt.Sprintf("cmbp-%d", i)

		// Call the createBucket function with the endpoint URL
		if err := createBucket(bucketName, endpointURL); err != nil {
			log.Println(err)
		}
	}
}
