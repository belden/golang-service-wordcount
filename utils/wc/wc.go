package Wc

import (
	"regexp"
	"strings"
)

type Dictionary struct {
	Total int
	Words map[string]int
}

func New() Dictionary {
	return Dictionary{0, map[string]int{}}
}

func SplitWords(Corpus string) []string {
	re := regexp.MustCompile("[A-Za-z'-]+")
	return re.FindAllString(Corpus, -1)
}

func CountWords(Words []string) Dictionary {
	Dict := New()

	for _, Word := range Words {
		Dict.Total++
		Dict.Words[strings.ToLower(Word)]++
	}

	return Dict
}

func Add(dicts ...Dictionary) Dictionary {
	Sum := New()

	for _, d := range dicts {
		Sum.Total += d.Total

		for k, v := range d.Words {
			if _, ok := Sum.Words[k]; ok {
				Sum.Words[k] += v
			} else {
				Sum.Words[k] = v
			}
		}
	}

	return Sum
}
