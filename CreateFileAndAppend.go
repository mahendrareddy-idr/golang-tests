package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	filePath := `C:\Users\Mahendra\Downloads\cdp_test_file.txt`

	// Open file in append mode, create if not exists
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Prepare content with timestamp (so every run is different)
	now := time.Now()
	content := fmt.Sprintf(
		"\n--- Update at %s ---\n"+
			"This file is used for CDP testing.\n"+
			"Each execution appends new content.\n"+
			"File size increases on every run.\n"+
			"Modification timestamp is updated.\n"+
			"Flush and close operations are performed properly.\n",
		now.Format("2006-01-02 15:04:05"),
	)

	// Append content
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Force flush to disk
	err = file.Sync()
	if err != nil {
		fmt.Println("Error syncing file:", err)
		return
	}

	// Explicitly update timestamps
	err = os.Chtimes(filePath, now, now)
	if err != nil {
		fmt.Println("Error updating timestamps:", err)
		return
	}

	fmt.Println("Content appended successfully.")
}
