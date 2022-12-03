package internal

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Detects CRLF line endings
	crlfRe = regexp.MustCompile("\r\n")
	// Detects LF line endings
	lfRe = regexp.MustCompile("\n")
)

func DetectNewline(input []byte) string {
	crlfs := crlfRe.FindAllIndex(input, -1)
	lfs := lfRe.FindAllIndex(input, -1)
	if len(crlfs) > len(lfs) {
		return "\r\n"
	}
	return "\n"
}

var (
	leadingSpacesRegex = regexp.MustCompile(`^[\s]+`)
)

func MinIndent(lines []string) int {
	minSpaces := 0
	for _, line := range lines {
		spaces := leadingSpacesRegex.FindStringIndex(line)
		if len(spaces) > 0 {
			minSpaces = Min(minSpaces, spaces[len(spaces)-1])
		}
	}
	return minSpaces
}

func PreviewLines(lines []string, length int) []string {
	totalLines := len(lines)
	previewLines := Min(length, totalLines-1)
	lines = lines[:previewLines]
	omitted := totalLines - previewLines
	if omitted > 0 {
		lines = append(lines, fmt.Sprintf("... %d lines omitted", omitted))
	}
	return lines
}

func Lines(source string, newline string) []string {
	return strings.Split(source, newline)
}

func LineNumber(source string, pos int, newline string) int {
	before := source[0:pos]
	lines := Lines(before, newline)
	return len(lines)
}
