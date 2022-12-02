package main

import (
	"github.com/sabhiram/go-gitignore"
)

type sourceMap map[string]bool

func (s sourceMap) ValidCount() int {
	valid := 0
	for _, ok := range s {
		if ok {
			valid += 1
		}
	}
	return valid
}

func (s sourceMap) Ignore(gi *ignore.GitIgnore) int {
	before := s.ValidCount()
	for source, _ := range s {
		if gi.MatchesPath(source) {
			s[source] = false
		}
	}
	after := s.ValidCount()
	return before - after
}

func (s sourceMap) Add(sources ...string) {
	for _, source := range sources {
		s[source] = true
	}
}

func (s sourceMap) Valid() []string {
	var valid []string
	for source, ok := range s {
		if ok {
			valid = append(valid, source)
		}
	}
	return valid
}
