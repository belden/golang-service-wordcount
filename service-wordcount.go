package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "regexp"
)

type WordCount struct {
	total int
	words map[string]int
}

func load_word_data(corpus []byte, rw http.ResponseWriter) *WordCount {
	dict := &WordCount{total: 0, words: map[string]int{}}
	dict.total = 7
	dict.words["hello"] = 1
	dict.words["world"] = 2
	out_json, _ := json.Marshal(dict)
	rw.Write(out_json)
	return dict
}

func wc_file(rw http.ResponseWriter, request *http.Request) {
	// request_json, _ := json.Marshal(request)
	// rw.Write(request_json)

	file, _, err := request.FormFile("file")
	if err == nil {
		content, _ := ioutil.ReadAll(file)
		fmt.Printf("got data: %s", content)
		load_word_data(content, rw)
		// out_json, _ := json.Marshal(out)
		// rw.Write(out_json)
	} else {
		fmt.Fprintf(rw, "encoutered error: %s", err)
	}

	fmt.Fprintln(rw, "Success!")
}

func main() {
	http.HandleFunc("/", wc_file)
	http.ListenAndServe(":3000", nil)
}
