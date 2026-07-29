package main

import (
	"math/rand"
	"os"
	"time"
)

func main() {
	filename := "testfile.txt"
	f, _ := os.OpenFile(filename, os.O_RDWR, 0644)
	defer f.Close()

	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, 1)

	for i := 0; i < 3; i++ { // corrupt 3 bytes
		offset := rand.Int63n(64) // within first 64 bytes
		f.Seek(offset, 0)
		f.Read(buf)
		buf[0] ^= 0xFF // invert bits
		f.Seek(offset, 0)
		f.Write(buf)
	}
}
