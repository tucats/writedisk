package main

import (
	"fmt"
	"strconv"
	"strings"
)

func formatSize(size int) string {
	if size >= gigabytes {
		return fmt.Sprintf("%.2f GB", float64(size)/gigabytes)
	} else if size >= megabytes {
		return fmt.Sprintf("%.2f MB", float64(size)/megabytes)
	} else if size >= kilobytes {
		return fmt.Sprintf("%.2f KB", float64(size)/kilobytes)
	} else {
		return fmt.Sprintf("%d bytes", size)
	}
}

func parseSize(s string) (int, error) {
	scale := 1

	text := strings.TrimSpace(strings.ToLower(s))

	switch {
	case strings.HasSuffix(text, "mb"):
		scale = megabytes
		text = strings.TrimSuffix(text, "mb")

	case strings.HasSuffix(text, "m"):
		scale = megabytes
		text = strings.TrimSuffix(text, "m")

	case strings.HasSuffix(text, "gb"):
		scale = gigabytes
		text = strings.TrimSuffix(text, "gb")

	case strings.HasSuffix(text, "g"):
		scale = gigabytes
		text = strings.TrimSuffix(text, "g")

	case strings.HasSuffix(text, "kb"):
		scale = kilobytes
		text = strings.TrimSuffix(text, "kb")

	case strings.HasSuffix(text, "k"):
		scale = kilobytes
		text = strings.TrimSuffix(text, "k")
	}

	size, err := strconv.Atoi(text)

	return size * scale, err
}
