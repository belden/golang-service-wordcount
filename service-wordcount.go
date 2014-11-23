package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func split_words(corpus string) []string {
	words := regexp.MustCompile("[A-Za-z'-]+")
	return words.FindAllString(corpus, -1)
}

type Dictionary struct {
	total int
	words map[string]int
}

func count_words(words []string) Dictionary {

	dict := Dictionary{total: 0, words: map[string]int{}}

	for _, word := range words {
		dict.total++
		dict.words[word]++
	}

	return dict
}

func wc_file(rw http.ResponseWriter, request *http.Request) {
	// for now I'll just assume it's a POST
	file, _, err := request.FormFile("file")

	// and I'm assuming the file is within the 10MB bounds

	if err == nil {
		// slurp the file, then convert byte array to a string
		corpus_bytes, _ := ioutil.ReadAll(file)
		corpus := string(corpus_bytes[:])

		// split into an array of words, then count them
		counts := count_words(split_words(corpus))

		// let's see what we've got here
		fmt.Printf("total: %d\n", counts.total)
		for word, freq := range counts.words {
			fmt.Printf("  %s: %d\n", word, freq)
		}
	} else {
		fmt.Fprintf(rw, "encoutered error: %s", err)
	}

	fmt.Fprintln(rw, "Success!")
}

func main() {
	http.HandleFunc("/", wc_file)
	http.ListenAndServe(":3000", nil)
}
