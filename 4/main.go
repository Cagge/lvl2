package main

import (
	"fmt"
	"sort"
	"strings"
)

func sortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}
func findAnagrams(words []string) map[string][]string {
	anagramSets := make(map[string][]string)
	for _, word := range words {
		sortedWord := sortString(strings.ToLower(word))
		anagramSets[sortedWord] = append(anagramSets[sortedWord], strings.ToLower(word))
	}
	result := make(map[string][]string)
	for _, value := range anagramSets {
		if len(value) > 1 {
			key := value[0]
			sort.Strings(value)
			result[key] = value
		}
	}

	return result
}

func main() {
	list := []string{
		"столик",
		"пятак",
		"тяпка",
		"алгоритм",
		"листок",
		"пятка",
		"слиток",
	}
	fmt.Println(findAnagrams(list))
}
