package main

import (
	"fmt"
)

func main() {
	multid := make([][]int, 0)
	multid = append(multid, []int{1, 2, 3})
	multid = append(multid, []int{4, 5, 6})

	for _, o := range multid {
		for _, o := range o {
			fmt.Printf("%d ", o)
		}
		fmt.Printf("\n")
	}
}
