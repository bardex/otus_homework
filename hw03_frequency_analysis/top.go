package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const punct = `'"!;:()-.,?`

func Top10(s string) []string {
	words := strings.Fields(s)
	wstats := make(map[string]int)
	for _, w := range words {
		w = strings.ToLower(strings.Trim(w, punct))
		if w != "" {
			wstats[w]++
		}
	}

	counts := []int{}
	wordsByCount := make(map[int][]string)
	for word, count := range wstats {
		if _, exists := wordsByCount[count]; !exists {
			counts = append(counts, count)
		}
		wordsByCount[count] = append(wordsByCount[count], word)
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})

	result := []string{}
	for _, count := range counts {
		words = wordsByCount[count]
		sort.Slice(words, func(i, j int) bool {
			return words[i] < words[j]
		})

		if len(result)+len(words) > 10 {
			tail := 10 - len(result)
			result = append(result, words[0:tail]...)
			break
		}

		result = append(result, words...)
	}

	return result
}
