package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	outputDir  = `C:\Users\Mahendra\Documents\GO-LANG\Dataset\LogsandSystemfiles`
	totalFiles = 600
)

var extensions = []string{
	".cs", ".cpp", ".py", ".js",
	".log",
	".ini", ".yaml", ".env",
}

func main() {
	rand.Seed(time.Now().UnixNano())
	os.MkdirAll(outputDir, os.ModePerm)

	for i := 1; i <= totalFiles; i++ {
		ext := extensions[i%len(extensions)]
		filename := fmt.Sprintf("file_%03d%s", i, ext)
		content := generateContent(ext, i)

		path := filepath.Join(outputDir, filename)
		os.WriteFile(path, []byte(content), 0644)
	}

	fmt.Println("✅ 100 source, log, and config files created successfully")
}

// ----------------------------------------------------

func generateContent(ext string, index int) string {
	switch ext {

	case ".cs":
		return fmt.Sprintf(`using System;

namespace SampleApp%d
{
    class Program
    {
        static void Main()
        {
            Console.WriteLine("Hello from C# file %d");
        }
    }
}`, index, index)

	case ".cpp":
		return fmt.Sprintf(`#include <iostream>
using namespace std;

int main() {
    cout << "Hello from C++ file %d" << endl;
    return 0;
}`, index)

	case ".py":
		return fmt.Sprintf(`def main():
    print("Hello from Python file %d")

if __name__ == "__main__":
    main()
`, index)

	case ".js":
		return fmt.Sprintf(`function main() {
    console.log("Hello from JavaScript file %d");
}

main();`, index)

	case ".log":
		return fmt.Sprintf(
			`[%s] INFO  Application started
[%s] DEBUG Processing request id=%d
[%s] WARN  Low memory warning
[%s] INFO  Application finished`,
			time.Now().Format(time.RFC3339),
			time.Now().Format(time.RFC3339),
			rand.Intn(10000),
			time.Now().Format(time.RFC3339),
			time.Now().Format(time.RFC3339),
		)

	case ".ini":
		return fmt.Sprintf(`[server]
port=808%d
host=localhost

[database]
user=admin
timeout=30`, index%10)

	case ".yaml":
		return fmt.Sprintf(`app:
  name: sample_app_%d
  version: 1.%d.%d
server:
  port: %d
  debug: true`,
			index,
			index%10,
			rand.Intn(10),
			8000+index,
		)

	case ".env":
		return fmt.Sprintf(`APP_NAME=sample_app_%d
APP_ENV=development
APP_PORT=%d
DEBUG=true`,
			index,
			3000+index,
		)
	}

	return ""
}
