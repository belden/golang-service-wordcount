package WcPaginate

import (
	"."
	"fmt"
	"testing"
)

func is_deeply(T *testing.T, got map[int]int, exp map[int]int) {
	for i := 0; i < len(exp); i++ {
		if got[i] != exp[i] {
			T.Error(fmt.Sprintf("mismatch at [%d]: got %d, expected %d from %s and %s\n", i, got[i], exp[i], got, exp))
		}
	}

	if len(got) != len(exp) {
		T.Error(fmt.Sprintf("size mismatch: got %d, expected %d\n", len(got), len(exp)))
	}
}

func is(T *testing.T, got string, exp string, message string) {
	if got != exp {
		T.Error(fmt.Sprintf("%s: got %s, exp %s\n", message, got, exp))
	}
}

func dumper(target map[int]int, message string) {
	for k, v := range target {
		fmt.Printf("%s: %d -> %d\n", message, k, v)
	}
}

func TestMakePageSlices(T *testing.T) {
	exp := map[int]int{
		0: 3,  // 1, 2, 3
		3: 6,  // 4, 5, 6
		6: 10, // 7, 8, 9, 10
	}
	got, err := WcPaginate.MakePageSlices(10, len(exp))
	is(T, fmt.Sprintf("%s", err), fmt.Sprintf("%s", nil), "no errors")
	is_deeply(T, got, exp)

	exp = map[int]int{
		0: 2,  // 1, 2
		2: 4,  // 3, 4
		4: 6,  // 5, 6
		6: 8,  // 7, 8
		8: 10, // 9, 10
	}
	got, err = WcPaginate.MakePageSlices(10, len(exp))
	is(T, fmt.Sprintf("%s", err), fmt.Sprintf("%s", nil), "no errors")
	is_deeply(T, got, exp)

	exp = map[int]int{
		0: 3,  // 1, 2, 3
		3: 6,  // 4, 5, 6
		6: 11, // 7, 8, 9, a, b
	}
	got, err = WcPaginate.MakePageSlices(11, len(exp))
	is(T, fmt.Sprintf("%s", err), fmt.Sprintf("%s", nil), "no errors")
	is_deeply(T, got, exp)

	exp = map[int]int{
		0: 4,  // 1, 2, 3, 4
		4: 8,  // 5, 6, 7, 8
		8: 12, // 9, a, b, c
	}
	got, err = WcPaginate.MakePageSlices(12, 3)
	is(T, fmt.Sprintf("%s", err), fmt.Sprintf("%s", nil), "no errors")
	is_deeply(T, got, exp)

	exp = map[int]int{
		0: 4,  // 1, 2, 3, 4
		4: 8,  // 5, 6, 7, 8
		8: 12, // 9, a, b, c
	}
	got, err = WcPaginate.MakePageSlices(0, len(exp))
	is(T, fmt.Sprintf("%s", err), "WcPaginate.MakePageSlices(page_data, pages): page_data must be > 0", "got an error")
}
