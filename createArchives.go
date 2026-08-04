package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	sourceDir  = `C:\Users\Mahendra\Downloads`
	outputDir  = `C:\Users\Mahendra\Documents\GO-LANG\Dataset\Pre-CompressedData`
	totalFiles = 100
)

var allowedExt = map[string]bool{
	".jpg": true,
	".csv": true,
}

func main() {
	os.MkdirAll(outputDir, os.ModePerm)

	files := getAllowedFiles()
	if len(files) == 0 {
		fmt.Println("❌ No .jpg or .csv files found in Downloads")
		return
	}

	// ZIP
	for i := 1; i <= totalFiles; i++ {
		createZip(files, fmt.Sprintf("data_%03d.zip", i))
	}

	// TAR.GZ
	for i := 1; i <= totalFiles; i++ {
		createTarGz(files, fmt.Sprintf("data_%03d.tar.gz", i))
	}

	// 7Z
	for i := 1; i <= totalFiles; i++ {
		create7z(files, fmt.Sprintf("data_%03d.7z", i))
	}

	// RAR
	for i := 1; i <= totalFiles; i++ {
		createRar(files, fmt.Sprintf("data_%03d.rar", i))
	}

	fmt.Println("✅ All archives created successfully")
}

// -----------------------------------------------------

func getAllowedFiles() []string {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		panic(err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if allowedExt[ext] {
			files = append(files, filepath.Join(sourceDir, e.Name()))
		}
	}
	return files
}

// -----------------------------------------------------
// ZIP

func createZip(files []string, name string) {
	out, err := os.Create(filepath.Join(outputDir, name))
	if err != nil {
		fmt.Println("ZIP create error:", err)
		return
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		info, _ := f.Stat()

		header, _ := zip.FileInfoHeader(info)
		header.Name = filepath.Base(file)

		w, _ := zw.CreateHeader(header)
		io.Copy(w, f)
		f.Close()
	}
}

// -----------------------------------------------------
// TAR.GZ

func createTarGz(files []string, name string) {
	out, err := os.Create(filepath.Join(outputDir, name))
	if err != nil {
		fmt.Println("TAR.GZ create error:", err)
		return
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		info, _ := f.Stat()

		header, _ := tar.FileInfoHeader(info, "")
		header.Name = filepath.Base(file)

		tw.WriteHeader(header)
		io.Copy(tw, f)
		f.Close()
	}
}

// -----------------------------------------------------
// 7Z (requires 7-Zip)

func create7z(files []string, name string) {
	args := []string{"a", filepath.Join(outputDir, name)}
	args = append(args, files...)

	cmd := exec.Command("7z", args...)
	// If 7z not in PATH, use:
	// cmd := exec.Command(`C:\Program Files\7-Zip\7z.exe`, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("7z error:", err)
	}
}

// -----------------------------------------------------
// RAR (requires WinRAR)

func createRar(files []string, name string) {
	args := []string{"a", filepath.Join(outputDir, name)}
	args = append(args, files...)

	cmd := exec.Command("rar", args...)
	// If rar not in PATH, use:
	// cmd := exec.Command(`C:\Program Files\WinRAR\rar.exe`, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("rar error:", err)
	}
}
