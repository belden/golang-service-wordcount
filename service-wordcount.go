package main

import (
	"./utils"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func FilenameInRequest(request *http.Request) (string, error) {
	request.ParseForm()
	// params := request.Form
	filename := request.FormValue("filename")
	if len(filename) > 0 {
		return filename, nil
	} else {
		return "", errors.New("no filename given in request")
	}
}

func wc_file() http.HandlerFunc {
	Cache := make(map[string]Wc.Dictionary)

	// for now I'll just assume it's a POST - assuming it's within 10MB, too
	return func(rw http.ResponseWriter, request *http.Request) {

		if request.Method == "GET" {

			if filename, err := FilenameInRequest(request); err != nil {

				// eg, http://localhost:3000/ - return the cached filenames
				filenames := make([]string, 0, len(Cache))
				for filename := range Cache {
					filenames = append(filenames, filename)
				}
				emit_json(rw, filenames)

			} else {

				// eg, http://localhost:3000/?filename=foo - return Cache["foo"]
				emit_json(rw, Cache[filename])

			}

		} else if request.Method == "DELETE" {

			filename, _ := FilenameInRequest(request)
			fmt.Printf("got a DELETE %s\n", filename)
			delete(Cache, filename)
			emit_json(rw, []byte(nil))

		} else if request.Method == "POST" {
			file, _, err := request.FormFile("file")

			filename, _ := FilenameInRequest(request)

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
				counts := Wc.CountWords(Wc.SplitWords(corpus))

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
	http.HandleFunc("/files", wc_file())
	http.ListenAndServe(portString, nil)
}
