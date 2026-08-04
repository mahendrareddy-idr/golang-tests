package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Folder structure
	baseDir := "Dataset"
	regularDir := filepath.Join(baseDir, "RegularFiles")

	// Create folders if they don't exist
	err := os.MkdirAll(regularDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}

	// File types and number of files
	fileTypes := []string{"txt", "csv", "rtf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf"}
	totalFiles := 100

	// Generate files
	for _, ext := range fileTypes {
		for i := 1; i <= totalFiles; i++ {
			filename := fmt.Sprintf("%s_%d.%s", ext, i, ext)
			filePath := filepath.Join(regularDir, filename)

			file, err := os.Create(filePath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				continue
			}

			// Write lines: 1 in first file, 2 in second, etc.
			for line := 1; line <= i; line++ {
				content := fmt.Sprintf("This is line %d for file %s\n", line, filename)

				// For RTF, DOC, DOCX, XLS, XLSX, PPT, PPTX, PDF, we just write plain text placeholders
				if ext == "rtf" {
					content = fmt.Sprintf("{\\rtf1 %s}", content)
				}

				_, err := file.WriteString(content)
				if err != nil {
					fmt.Println("Error writing to file:", err)
					break
				}
			}
			file.Close()
		}
	}

	fmt.Println("All files generated inside", regularDir)
}
