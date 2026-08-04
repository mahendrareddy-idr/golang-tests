package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	outputDir  = `C:\Users\Mahendra\Documents\GO-LANG\Dataset\LockedFiles`
	totalFiles = 100
	vhdxSizeMB = 10
)

func main() {
	os.MkdirAll(outputDir, os.ModePerm)

	// 1️⃣ Create REAL VHDX files
	for i := 1; i <= totalFiles; i++ {
		name := fmt.Sprintf("disk_%03d.vhdx", i)
		createVHDX(name)
	}

	// 2️⃣ Create LOCKED PST files
	for i := 1; i <= totalFiles; i++ {
		name := fmt.Sprintf("mail_%03d.pst", i)
		createLockedBinary(name)
	}

	// 3️⃣ Create LOCKED OST files
	for i := 1; i <= totalFiles; i++ {
		name := fmt.Sprintf("mail_%03d.ost", i)
		createLockedBinary(name)
	}

	fmt.Println("✅ VHDX created (real), PST/OST created (locked placeholders)")
}

// ----------------------------------------------------
// REAL VHDX creation (Windows-native)

func createVHDX(name string) {
	path := filepath.Join(outputDir, name)

	cmd := exec.Command(
		"powershell",
		"-Command",
		fmt.Sprintf(
			"New-VHD -Path '%s' -SizeBytes %dMB -Dynamic",
			path,
			vhdxSizeMB,
		),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// ----------------------------------------------------
// Locked PST / OST placeholders

func createLockedBinary(name string) {
	path := filepath.Join(outputDir, name)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	// Write binary-like data
	f.Write([]byte{
		0xD0, 0xCF, 0x11, 0xE0,
		0xA1, 0xB1, 0x1A, 0xE1,
	})

	// Hold lock briefly (simulate locked system file)
	go func(file *os.File) {
		time.Sleep(5 * time.Second)
		file.Close()
	}(f)
}
