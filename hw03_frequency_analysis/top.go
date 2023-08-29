package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type stat struct {
	word string
	cnt  int
}

func Top10(src string) []string {
	// слит строки в слова
	words := strings.Fields(src)
	// статистика
	cnt := make(map[string]int)
	for _, w := range words {
		cnt[w]++
	}

	// переделать в массив
	arr := make([]stat, 0, len(cnt))
	for k, v := range cnt {
		arr = append(arr, stat{word: k, cnt: v})
	}

	// отсортировать массив cnt desc word asc
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].cnt > arr[j].cnt ||
			(arr[i].cnt == arr[j].cnt && arr[i].word < arr[j].word)
	})

	// top 10
	result := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(arr); i++ {
		result = append(result, arr[i].word)
	}

	return result
}
