package WcPaginate

import (
	"errors"
)

func MakePageSlices(space int, pages int) (map[int]int, error) {
	if space < 1 {
		return nil, errors.New("WcPaginate.MakePageSlices(page_data, pages): page_data must be > 0")
	}

	tuples := make(map[int]int)

	// MakePageSlices(10, 3) -> {0:3, 3:6, 6:10}
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

	return tuples, nil
}
