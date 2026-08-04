package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	totalFiles   = 170_000_000
	workerCount  = 30     // adjust based on CPU/disk (100–5000 typical)
	batchSize    = 50_000 // controls memory pressure
	outputFolder = "./output"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range jobs {
		// Create filename (you can shard dirs for performance)
		dir := filepath.Join(outputFolder, fmt.Sprintf("%02d", i%1000))
		_ = os.MkdirAll(dir, os.ModePerm)

		filePath := filepath.Join(dir, strconv.Itoa(i))

		// Create 1-byte file
		f, err := os.Create(filePath)
		if err == nil {
			_, _ = f.Write([]byte{0})
			f.Close()
		}
	}
}

func main() {
	jobs := make(chan int, batchSize)

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// Feed jobs
	go func() {
		for i := 0; i < totalFiles; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Wait for completion
	wg.Wait()

	fmt.Println("Done creating files.")
}
