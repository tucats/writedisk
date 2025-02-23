package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func runThreads(path string, threadCount int, fileCount int, fileSize int, startValue int, logging bool) {
	var wg sync.WaitGroup

	start := time.Now()

	if logging {
		fmt.Printf("Write %d files of size %s bytes to %s\n", fileCount, formatSize(fileSize), path)
	}

	// Verify that the path exists and is writable
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)

		return
	}

	// Try to write a file to the path, and delete it. If this fails, bail out.
	if err := os.WriteFile(filepath.Join(path, fileBase+"-probe"+fileExtension), make([]byte, fileSize), 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)

		return
	} else {
		err = os.Remove(filepath.Join(path, fileBase+"-probe"+fileExtension))
		if err != nil {
			fmt.Printf("Error removing file: %v\n", err)
		}
	}

	fmt.Printf("Launching %d write operation threads...\n", threadCount)

	// Create a buffer with all possible byte values from the given start value.
	buffer := make([]byte, fileSize)
	for i := 0; i < fileSize; i++ {
		buffer[i] = byte((i + startValue) % 256)
	}

	// Launch one thread for each thread count, and pass it a range from the total
	// count.
	for thread := 0; thread < threadCount; thread++ {
		wg.Add(1)

		go worker(path, thread, thread*fileCount/threadCount, (thread+1)*fileCount/threadCount, buffer, logging, &wg)
	}

	if logging {
		fmt.Printf("Waiting for write operations to complete...\n")
	}

	wg.Wait()

	speed := float64(fileCount*fileSize) / time.Since(start).Seconds()
	elapsed := time.Since(start).String()

	fmt.Printf("Wrote %s in %s, %s/second\n", formatSize(fileCount*fileSize), elapsed, formatSize(int(speed)))
}

func worker(path string, thread, start, end int, buffer []byte, logging bool, wg *sync.WaitGroup) {
	defer wg.Done()

	n := 0
	useDuration := false

	each := 10
	if len(buffer) > 10*gigabytes {
		each = 1
	}

	now := time.Now()
	duration := time.Second * 10

	for j := start; j < end; j++ {
		n++

		if logging {
			log := false

			// If IT's been long enough since the last long, we switch to
			// duration logging mode.
			if time.Since(now) > duration {
				useDuration = true
			}

			// If we are in duration loggin gmode, we log each time the duration
			// has passed
			if useDuration && (time.Since(now) > duration) {
				log = true
			}

			// If we're are not in duraiton loggin gmode, we log each time the number of
			// files has been written is a multiple of 'each'.
			if !useDuration && n%each == 0 {
				log = true
			}

			// If after all that we've decided to log where we are, do that now. Always
			// record when we did the logging so we know if/when to switch to duration
			// logging.
			if log {
				fmt.Printf("Thread %3d: has written %4d files\n", thread+1, n)

				now = time.Now()
			}
		}

		filePath := filepath.Join(path, fmt.Sprintf("%s-%03d-%08d%s", fileBase, thread+1, j+1, fileExtension))

		if err := os.WriteFile(filePath, buffer, 0644); err != nil {
			fmt.Printf("Error writing file: %v\n", err)
		}
	}
}
