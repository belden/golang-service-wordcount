package Wc

import (
	"."
	"fmt"
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
