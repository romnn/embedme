package internal

import (
	"regexp"
)

// Match ...
type Match struct {
	Text  string
	Start int
	End   int
}

// GetMatches ...
func GetMatches(exp *regexp.Regexp, text string) []map[string]Match {
	var results []map[string]Match
	matches := exp.FindAllStringSubmatchIndex(text, -1)
	for _, pos := range matches {
		result := make(map[string]Match)
		for _, name := range exp.SubexpNames() {
			i := exp.SubexpIndex(name)
			if i < 0 {
				// no match
				continue
			}
			start := pos[i*2+0]
			end := pos[i*2+1]
			if start < 0 || end < 0 {
				// no match
				continue
			}
			result[name] = Match{
				Text:  text[start:end],
				Start: start,
				End:   end,
			}
		}
		results = append(results, result)
	}
	return results
}
