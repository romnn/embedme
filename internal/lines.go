package internal

import (
	"fmt"
	"golang.org/x/exp/constraints"
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

func min[T constraints.Ordered](a T, b T) T {
	if a > b {
		return b
	}
	return a
}

// func min[T constraints.Ordered](s []T) T {
// 	if len(s) == 0 {
// 		var zero T
// 		return zero
// 	}
// 	m := s[0]
// 	for _, v := range s {
// 		if m > v {
// 			m = v
// 		}
// 	}
// 	return m
// }

var (
	leadingSpacesRegex = regexp.MustCompile(`^[\s]+`)
)

func GetMinimumSpaces(lines []string) int {
	minSpaces := 0
	for _, line := range lines {
		spaces := leadingSpacesRegex.FindStringIndex(line)
		if len(spaces) > 0 {
			minSpaces = min(minSpaces, spaces[len(spaces)-1])
		}
		// fmt.Println(line)
		// fmt.Println(spaces)
	}
	return minSpaces
}

func PreviewLines(lines []string, length int) []string {
	// newline := DetectNewline([]byte(source))
	// lines := GetLines(source, newline)
	totalLines := len(lines)
	previewLines := min(length, totalLines-1)
	lines = lines[:previewLines]
	// preview := strings.Join(lines, newline)
	omitted := totalLines - previewLines
  lines = append(lines, fmt.Sprintf("\n... %d lines omitted", omitted))
  return lines
	// if omitted > 0 {
	// 	preview += fmt.Sprintf("\n... %d lines omitted", omitted)
	// }
	// return preview
}

func GetLines(source string, newline string) []string {
	return strings.Split(source, newline)
}

func GetLineNumber(source string, pos int, newline string) int {
	before := source[0:pos]
	// lines := newlineRe.Split(before, -1)
	// lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
	lines := GetLines(before, newline)
	return len(lines)
}
