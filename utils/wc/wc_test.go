package Wc

import (
	"."
	"fmt"
	"testing"
)

func is_deeply(T *testing.T, Got Wc.Dictionary, Exp Wc.Dictionary) {
	// Totals better match up
	if Got.Total != Exp.Total {
		T.Error(fmt.Sprintf("Total: Got %d, Expected %d\n", Got.Total, Exp.Total))
	}

	// And word frequencies better match, too
	for ExpKey := range Exp.Words {
		if Got.Words[ExpKey] != Exp.Words[ExpKey] {
			T.Error(fmt.Sprintf("Words[%s]: Got %d, Expected %d\n", ExpKey, Got.Words[ExpKey], Exp.Words[ExpKey]))
		}
	}
}

func TestNew(T *testing.T) {
	Exp := Wc.Dictionary{Total: 0, Words: map[string]int{}}
	is_deeply(T, Wc.New(), Exp)
}

func TestSplitWords(T *testing.T) {
	Got := Wc.SplitWords("hello world this is a meaningful test, Hello hello")
	Exp := []string{"hello", "world", "this", "is", "a", "meaningful", "test", "Hello", "hello"}
	for i := range Got {
		if Got[i] != Exp[i] {
			T.Error(fmt.Sprintf("got %s, expected %s\n", Got[i], Exp[i]))
		}
	}
}

func TestCountWords(T *testing.T) {
	Got := Wc.CountWords([]string{"hello", "world", "HELLO", "moon"})
	Exp := Wc.Dictionary{Total: 4, Words: map[string]int{
		"hello": 2,
		"world": 1,
		"moon":  1,
	}}

	is_deeply(T, Got, Exp)
}

func TestAdd(T *testing.T) {
	corpus1 := "hello world this is a meaningful test, Hello hello"
	corpus2 := "goodnight moon, this isn't the end of the world"
	corpus3 := "moon over my hammy, bacon eggs and cheese"
	Dict1 := Wc.CountWords(Wc.SplitWords(corpus1))
	Dict2 := Wc.CountWords(Wc.SplitWords(corpus2))
	Dict3 := Wc.CountWords(Wc.SplitWords(corpus3))

	entirety := fmt.Sprintf("%s\n%s\n%s\n", corpus1, corpus2, corpus3)
	Exp := Wc.CountWords(Wc.SplitWords(entirety))

	// this is the line of code under test
	Got := Wc.Add(Dict1, Dict2, Dict3)

	is_deeply(T, Got, Exp)
}
