package stringmatch

import (
	"sort"

	aca "github.com/BobuSumisu/aho-corasick"
)

func MatchString(str string, targets []string) []string {
	trie := aca.NewTrieBuilder().
		AddStrings(targets).
		Build()

	matches := trie.MatchString(str)
	sortMatches(matches)

	var (
		results         []string
		minAvailablePos int
	)

	for i := range matches {
		match := matches[i]

		pos := int(match.Pos())

		if pos < minAvailablePos {
			continue
		}

		minAvailablePos = pos + len(match.Match())
		results = append(results, match.MatchString())
	}

	return results
}

// sort by following priority:
// 1. index of first matched alphabet
// 2. length of pattern
// 3. dictionary order
func sortMatches(matches []*aca.Match) {
	// sort by dictionary
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].String() < matches[j].String()
	})

	// sort by len
	sort.SliceStable(matches, func(i, j int) bool {
		return len(matches[i].Match()) > len(matches[j].Match())
	})

	// sort by pos
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].Pos() < matches[j].Pos()
	})
}
