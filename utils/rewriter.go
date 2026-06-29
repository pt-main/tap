package utils

import (
	"fmt"
	"strings"
)

// Rewriter provides in-place terminal output rewriting using ANSI escape codes.
type Rewriter struct {
	lastLines int
}

func NewRewriter() *Rewriter {
	return &Rewriter{}
}

// Write replaces the previously written text with new content without scrolling.
func (r *Rewriter) Write(text string) {
	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	currentLines := len(lines) - 1
	if r.lastLines > 0 {
		fmt.Printf("\033[%dA", r.lastLines)
	}
	for i, line := range lines {
		fmt.Print("\033[2K\r")
		if i < len(lines)-1 {
			fmt.Println(line)
		} else {
			fmt.Print(line)
		}
	}
	if currentLines < r.lastLines {
		extra := r.lastLines - currentLines
		for i := 0; i < extra; i++ {
			fmt.Print("\033[1B")
			fmt.Print("\033[2K\r")
		}

		fmt.Printf("\033[%dA", extra)
	}
	r.lastLines = currentLines
}
