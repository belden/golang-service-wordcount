package Wc

import (
	"regexp"
	"strings"
)

type Dictionary struct {
	Total int
	Words map[string]int
}

func SplitWords(Corpus string) []string {
	re := regexp.MustCompile("[A-Za-z'-]+")
	return re.FindAllString(Corpus, -1)
}

func CountWords(Words []string) Dictionary {
	Dict := Dictionary{Total: 0, Words: map[string]int{}}

	for _, Word := range Words {
		Dict.Total++
		Dict.Words[strings.ToLower(Word)]++
	}

	return Dict
}
