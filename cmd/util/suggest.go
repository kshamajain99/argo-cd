package util

import (
	"fmt"
	"strings"
)

const (
	minSuggestionsDistance = 3
)

var (
	argocdCommands = []string{"argocd-util", "argocd-server", "argocd-application-controller", "argocd-repo-server", "argocd-dex"}
)

func FindSuggestions(binaryName string) string {
	suggestionsString := ""
	if suggestions := SuggestionsFor(binaryName); len(suggestions) > 0 {
		suggestionsString += "\n\nDid you mean this?\n"
		for _, s := range suggestions {
			suggestionsString += fmt.Sprintf("\t%v\n", s)
		}
	}
	return suggestionsString
}

// Logic from https://github.com/spf13/cobra/blob/master/command.go#L721
// Modified suggestion by prefix to suggest when either binary Name is
// prefix to command name or command name is prefix to binary name
// Added suggestions by substring
// SuggestionsFor provides suggestions for the binary name.
func SuggestionsFor(binaryName string) []string {
	var suggestions []string
	for _, cmdName := range argocdCommands {
		levenshteinDistance := ld(binaryName, cmdName, true)
		suggestByLevenshtein := levenshteinDistance <= minSuggestionsDistance
		suggestByPrefix := strings.HasPrefix(strings.ToLower(cmdName), strings.ToLower(binaryName)) ||
			strings.HasPrefix(strings.ToLower(binaryName), strings.ToLower(cmdName))
		suggestBySubstring := strings.Contains(cmdName, binaryName)
		if suggestByLevenshtein || suggestByPrefix || suggestBySubstring {
			suggestions = append(suggestions, cmdName)
		}
	}
	return suggestions
}

// Logic from https://github.com/spf13/cobra/blob/master/cobra.go#L165
// ld compares two strings and returns the levenshtein distance between them.
func ld(s, t string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}
