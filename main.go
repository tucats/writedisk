package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	kilobytes = 1024
	megabytes = 1024 * kilobytes
	gigabytes = 1024 * megabytes
)

var fileBase = "file-" + uuid.New().String()[:8]

func main() {
	var (
		path    string
		count   int
		size    int
		logging bool
		err     error
	)

	count = 1
	size = 10 * megabytes

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "--count", "-c":
			count, err = strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Printf("Invalid count: %s\n", os.Args[i+1])
				os.Exit(1)
			}

			i++

		case "--size", "-s":
			size, err = parseSize(os.Args[i+1])
			if err != nil {
				fmt.Printf("Invalid size: %s\n", os.Args[i+1])
				os.Exit(1)
			}

			i++

		case "--logging", "-l", "-v", "--verbose":
			logging = true

		case "--path", "-p":
			path = os.Args[i+1]
			i++

		case "--help", "-h":
			fmt.Println("writedisk creates a specified number of files with a given size, filled with non-zero data.")
			fmt.Println()
			fmt.Println("Usage: writedisk [options]")
			fmt.Println("\nOptions:")
			fmt.Println("  --count, -c <number>: Number of files to create (default: 1)")
			fmt.Println("  --size, -s <size>   : Size of each file in bytes (default: 10MB)")
			fmt.Println("  --verbose, -v       : Enable logging")
			fmt.Println("  --path, -p <path>   : Output path for the files (required)")
			os.Exit(0)

		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Printf("Invalid argument: %s\n", arg)
				os.Exit(1)
			}

			path = os.Args[i]
		}
	}

	if path == "" {
		fmt.Println("No output path specified")
		os.Exit(1)
	}

	fmt.Printf("Writing %d files, for a total size of %s\n", count, formatSize(count*size))

	launch(path, count, size, logging)
}

func launch(path string, count int, size int, logging bool) {
	var wg sync.WaitGroup

	start := time.Now()

	threadCount := runtime.NumCPU() * 2

	if logging {
		fmt.Printf("Using %d threads to write files\n", threadCount)
		fmt.Printf("Write %d files of size %s bytes to %s\n", count, formatSize(size), path)
	}

	// Verify that the path exists and is writable
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)

		return
	}

	// Try to write a file to the path, and delete it. If this fails, bail out.
	if err := os.WriteFile(filepath.Join(path, fileBase+"-probe.txt"), make([]byte, size), 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)

		return
	} else {
		err = os.Remove(filepath.Join(path, fileBase+"-probe.txt"))
		if err != nil {
			fmt.Printf("Error removing file: %v\n", err)
		}
	}

	fmt.Printf("launching write operations...\n")

	// Create a buffer with all possible byte values. The buffer simply increments,
	// but starts with a random value so two files generated with different commands
	// are unlikely to have the same content.
	startValue := 0
	if start, err := rand.Int(rand.Reader, big.NewInt(256)); err != nil {
		startValue = int(start.Int64())
	}

	buffer := make([]byte, size)
	for i := 0; i < size; i++ {
		buffer[i] = byte((i + startValue) % 256)
	}

	// Launch one thread for each thread count, and pass it a range from the total
	// count.
	for thread := 0; thread < threadCount; thread++ {
		wg.Add(1)

		go func(i, start, end int, buffer []byte, logging bool) {
			defer wg.Done()

			n := 0

			each := 10
			if size > 10*megabytes {
				each = 1
			}

			now := time.Now()
			duration := time.Second * 10

			for j := start; j < end; j++ {
				n++
				if logging && (time.Since(now) > duration || n%each == 0) {
					fmt.Printf("Thread %3d: wrote %4d files\n", i, n)

					now = time.Now()
				}

				filePath := filepath.Join(path, fmt.Sprintf("%s-%03d-%08d.txt", fileBase, thread, j))

				if err := os.WriteFile(filePath, buffer, 0644); err != nil {
					fmt.Printf("Error writing file: %v\n", err)
				}
			}
		}(thread, thread*count/threadCount, (thread+1)*count/threadCount, buffer, logging)
	}

	if logging {
		fmt.Printf("Waiting for write operations to complete...\n")
	}

	wg.Wait()
	fmt.Printf("Finished write operations in %v\n", time.Since(start))
}
