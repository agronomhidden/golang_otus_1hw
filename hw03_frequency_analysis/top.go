package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

var ErrIncorrectWord = errors.New("input is not a correct word")
var re = regexp.MustCompile(`[^а-яa-z-]`)
var separator = regexp.MustCompile(`\s`)
var top = 10

type frequencies struct {
	word      string
	frequency int
}

type frequencyContainer struct {
	index     map[string]*frequencies
	container []*frequencies
}

func (e *frequencyContainer) AppendWord(word string) {
	if item, ok := e.index[word]; ok {
		item.frequency++
	} else {
		if e.index == nil {
			e.index = make(map[string]*frequencies)
		}
		freq := &frequencies{word, 1} // можно использовать zero-value,
		e.index[word] = freq          // но для точности пусть будут фактические частоты
		e.container = append(e.container, freq)
	}
}

func (e frequencyContainer) GetTopWords(top int) (result []string) {
	sort.Slice(e.container, func(i, j int) bool {
		a := e.container[i]
		b := e.container[j]
		if a.frequency == b.frequency { // для более ожидаемого результата делаем сравнение слов,
			return a.word < b.word // если частоты одинаковы
		}
		return a.frequency > b.frequency
	})
	for i, item := range e.container {
		if i == top {
			break
		}
		result = append(result, item.word)
	}
	return
}

func Top10(input string) []string {
	container := frequencyContainer{}
	words := separator.Split(input, -1)
	for _, word := range words {
		nWord, err := normalize(word)
		if err != nil {
			continue
		}
		container.AppendWord(nWord)
	}
	return container.GetTopWords(top)
}

func normalize(word string) (string, error) {
	word = re.ReplaceAllString(strings.ToLower(word), "")
	if len(word) == 0 || word == "-" {
		return "", ErrIncorrectWord
	}
	return word, nil
}
