package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wc struct {
	word  string
	count int
}

func Top10(s string) []string {
	words := strings.Fields(s)

	wcMap := make(map[string]int)
	for _, word := range words {
		if count, ok := wcMap[word]; ok {
			wcMap[word] = count + 1
		} else {
			wcMap[word] = 1
		}
	}

	i := 0
	wcSlice := make([]wc, len(wcMap))
	for word, count := range wcMap {
		wcSlice[i] = wc{word, count}
		i++
	}

	sort.Slice(wcSlice, func(i, j int) bool {
		if wcSlice[i].count == wcSlice[j].count {
			return wcSlice[i].word < wcSlice[j].word
		}
		return wcSlice[i].count > wcSlice[j].count
	})

	length := 10
	if len(wcSlice) < 10 {
		length = len(wcSlice)
	}
	top10 := make([]string, length)
	for i := 0; i < length; i++ {
		top10[i] = wcSlice[i].word
	}

	return top10
}
