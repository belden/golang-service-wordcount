package main

import (
	"encoding/json"
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
	Total int
	Words map[string]int
}

func emit_json(rw http.ResponseWriter, target Dictionary) {
	js, err := json.Marshal(target)

	rw.Header().Set("Content-Type", "application/json")
	if err == nil {
		rw.WriteHeader(200)
		rw.Write(js)
	} else {
		rw.WriteHeader(500)
		rw.Write([]byte(`"{\"error\": \"yikes\"}"`))
	}
}

func count_words(words []string) Dictionary {

	dict := Dictionary{Total: 0, Words: map[string]int{}}

	for _, word := range words {
		dict.Total++
		dict.Words[word]++
	}

	return dict
}

func wc_file(cache map[string]Dictionary) http.HandlerFunc {
	// for now I'll just assume it's a POST - assuming it's within 10MB, too
	return func(rw http.ResponseWriter, request *http.Request) {
		file, _, err := request.FormFile("file")

		// grab the filename
		request.ParseForm()
		params := request.Form
		fn := params["filename"]
		filename := fn[0]

		// bail out if it's in cache already
		if seen, ok := cache[filename]; ok {
			fmt.Printf("Returning cached data for %s\n", filename)
			emit_json(rw, seen)
			return
		}

		if err == nil {
			// slurp the file, then convert byte array to a string
			corpus_bytes, _ := ioutil.ReadAll(file)
			corpus := string(corpus_bytes[:])

			// split into an array of words, then count them
			counts := count_words(split_words(corpus))

			// store in cache
			fmt.Printf("storing data for %s in cache\n", filename)
			cache[filename] = counts

			// send response
			emit_json(rw, counts)
		} else {
			fmt.Fprintf(rw, "encountered error: %s", err)
		}
	}
}

func main() {
	cache := make(map[string]Dictionary)
	http.HandleFunc("/", wc_file(cache))
	http.ListenAndServe(":3000", nil)
}
