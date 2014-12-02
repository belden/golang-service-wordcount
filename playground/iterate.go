package main

import (
	"fmt"
)

func GetFoo(Size int) []int {
	var Out []int
	for i := 1; i < Size+1; i++ {
		Out = append(Out, i)
	}
	return Out
}

func MakeOffsets(space int, pages int) map[int]int {
	tuples := make(map[int]int)

	// MakeOffsets(10, 3) -> {0:3, 3:6, 6:10}
	// note we handle widowing on the last page
	widow := space % pages
	avg_page_size := (space - widow) / pages
	last_page_size := avg_page_size + widow

	_i := 0
	for i := 0; i <= (space - last_page_size); i += avg_page_size {
		tuples[i] = i + avg_page_size
		_i = i
	}
	tuples[_i] = space

	return tuples
}

func main() {
	Foo := GetFoo(10)

	tuples := MakeOffsets(len(Foo), 3)

	pagesC := make(chan []int)

	for low, high := range tuples {
		low, high := low, high
		go func() {
			page := Foo[low:high]
			pagesC <- page
		}()
	}

	var pages []interface{} // the old "I don't know how to make a multidimensional array"
	for i := 0; i < len(tuples); i++ {
		it := <-pagesC
		pages = append(pages, it)
		fmt.Printf("%d\n", it)
	}

}
