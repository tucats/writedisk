package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	kilobytes     = 1024
	megabytes     = 1024 * kilobytes
	gigabytes     = 1024 * megabytes
	fileExtension = ".data"
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
	threadCount := runtime.NumCPU() * 2
	maxThreads := runtime.NumCPU() * 10

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "--count", "-c":
			count, err = strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Printf("Invalid count: %s\n", os.Args[i+1])
				os.Exit(1)
			}

			if count < 1 {
				fmt.Printf("Invalid count: %d, minimum value is 1\n", count)
				os.Exit(1)
			}

			i++

		case "--threads", "--thread", "-t":
			threadCount, err = strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Printf("Invalid threads: %s\n", os.Args[i+1])
				os.Exit(1)
			}

			if threadCount < 1 {
				fmt.Printf("Invalid threads: %d, minimum value is 1\n", threadCount)
				os.Exit(1)
			}

			if threadCount > maxThreads {
				fmt.Printf("Invalid threads: %d, maximum value is %d\n", threadCount, maxThreads)
				os.Exit(1)
			}

			i++

		case "--size", "-s":
			size, err = parseSize(os.Args[i+1])
			if err != nil {
				fmt.Printf("Invalid size: %s\n", os.Args[i+1])
				os.Exit(1)
			}

			if size < 1 {
				fmt.Printf("Invalid size: %d, minimum value is 1\n", size)
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
			fmt.Println("  --count, -c <number>:   Number of files to create (default: 1)")
			fmt.Println("  --path, -p <path>   :   Output path for the files (required)")
			fmt.Println("  --size, -s <size>   :   Size of each file in bytes (default: 10MB)")
			fmt.Println("  --threads, -s <num> :   Number of threads to use (default: twice the number of CPU cores)")
			fmt.Println("  --verbose, -v       :   Enable logging")
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

	fileBase = fileBase + "-" + formatSize(size)

	if threadCount > count {
		threadCount = count
	}

	fmt.Printf("Writing %d files, for a total size of %s\n", count, formatSize(count*size))

	runThreads(path, threadCount, count, size, logging)
}
