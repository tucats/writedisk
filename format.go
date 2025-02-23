package main

import (
	"fmt"
	"strconv"
	"strings"
)

func formatSize(size int) string {
	var text, scale string

	if size >= gigabytes {
		text = fmt.Sprintf("%.2f", float64(size)/gigabytes)
		scale = "GB"
	} else if size >= megabytes {
		text = fmt.Sprintf("%.2f", float64(size)/megabytes)
		scale = "MB"
	} else if size >= kilobytes {
		text = fmt.Sprintf("%.2f", float64(size)/kilobytes)
		scale = "KB"
	} else {
		text = fmt.Sprintf("%d", size)
		scale = "B"
	}

	if strings.HasSuffix(text, ".00") {
		text = strings.TrimSuffix(text, ".00")
	}

	return text + scale
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
