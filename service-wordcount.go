package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func split_words(corpus string) []string {
	words := regexp.MustCompile("[A-Za-z'-]+")
	return words.FindAllString(corpus, -1)
}

type Dictionary struct {
	Total int
	Words map[string]int
}

func emit_json(rw http.ResponseWriter, target interface{}) {
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
		dict.Words[strings.ToLower(word)]++
	}

	return dict
}

func filename(request *http.Request) string {
	request.ParseForm()
	params := request.Form
	fn := params["filename"]
	return fn[0]
}

func wc_file() http.HandlerFunc {
	Cache := make(map[string]Dictionary)

	// for now I'll just assume it's a POST - assuming it's within 10MB, too
	return func(rw http.ResponseWriter, request *http.Request) {

		if request.Method == "GET" {

			filenames := make([]string, 0, len(Cache))
			for filename := range Cache {
				filenames = append(filenames, filename)
			}
			emit_json(rw, filenames)

		} else if request.Method == "DELETE" {

			filename := filename(request)
			fmt.Printf("got a DELETE %s\n", filename)
			delete(Cache, filename)
			emit_json(rw, []byte(nil))

		} else if request.Method == "POST" {
			file, _, err := request.FormFile("file")
			filename := filename(request)

			// bail out if it's in Cache already
			if seen, ok := Cache[filename]; ok {
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

				// store in Cache
				fmt.Printf("storing data for %s in Cache\n", filename)
				Cache[filename] = counts

				// send response
				emit_json(rw, counts)
			} else {
				fmt.Fprintf(rw, "encountered error: %s", err)
			}
		}
	}
}

func main() {
	// read --port command-line option
	portNumber := flag.Int("port", 3000, "port to start on")
	flag.Parse()

	// tell the user where to find the service
	portString := fmt.Sprintf(":%d", *portNumber)
	fmt.Printf("Starting on http://localhost%s\n", portString)

	// register handlers and start listening for requests
	http.HandleFunc("/", wc_file())
	http.ListenAndServe(portString, nil)
}
