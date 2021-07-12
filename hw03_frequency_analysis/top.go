package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const punct = `'"!;:()-.,?`

func Top10(s string) []string {
	// разбиваем на слова по пробелу
	words := strings.Fields(s)
	wstats := make(map[string]int)
	for _, w := range words {
		// очищаем от знаков пунктуации и приводим к нижнему регистру
		w = strings.ToLower(strings.Trim(w, punct))
		// для не пустых слов прибавляем счетчик
		if w != "" {
			wstats[w]++
		}
	}

	// список полученных повторов слов, нужен для сортировки по популярности
	counts := []int{}
	// группировка: кол-во повторов и все слова с таким количеством, нужна для сортировки по-алфавиту
	wordsByCount := make(map[int][]string)
	for word, count := range wstats {
		if _, exists := wordsByCount[count]; !exists {
			counts = append(counts, count)
		}
		wordsByCount[count] = append(wordsByCount[count], word)
	}

	// сортируем список повторов слов по-убыванию
	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})

	result := []string{}
	for _, count := range counts {
		// все слова с данным количеством повторов
		words = wordsByCount[count]
		// сортируем их по-алфавиту
		sort.Slice(words, func(i, j int) bool {
			return words[i] < words[j]
		})
		// если после прибавления всех слов данной группы будет превышен лимит, то выбираем только нужное кол-во слов
		if len(result)+len(words) > 10 {
			tail := 10 - len(result)
			result = append(result, words[0:tail]...)
			break
		}
		// если нет - просто добавляем слова группы к результату
		result = append(result, words...)
	}

	return result
}
