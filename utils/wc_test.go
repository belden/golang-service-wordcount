package Wc

import (
	"."
	"fmt"
	"reflect"
	"testing"
)

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

	// make sure Wc.CountWords hands back a Wc.Dictionary
	GotType, ExpType := reflect.TypeOf(Got), reflect.TypeOf(Exp)
	if GotType != ExpType {
		T.Error(fmt.Sprintf("type mismatch: %s != %s", GotType, ExpType))
	}

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
