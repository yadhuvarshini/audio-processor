package utils

import (
	"strings"
)

func FakeExtraction(transcript string) []string {
	// For demo purposes, split the transcript into words and pick top 3
	words := strings.Fields(transcript)
	if len(words) > 3 {
		return words[:3]
	}
	return words
}
