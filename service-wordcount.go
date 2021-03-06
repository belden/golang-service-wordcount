package main

import (
	"./utils/wc"
	"./utils/wc/paginate"
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
		rw.Write([]byte(`{"error": "error in json encoding"}`))
	}
}

func FilenamesInRequest(request *http.Request) ([]string, error) {
	request.ParseForm()
	params := request.Form

	var filenames []string
	for _, filename := range params["filename"] {
		filenames = append(filenames, filename)
	}

	if len(filenames) > 0 {
		return filenames, nil
	} else {
		return nil, errors.New("no filenames given in request")
	}
}

func InitCache() map[string]Wc.Dictionary {
	return make(map[string]Wc.Dictionary)
}

func HandleWcFile(WcCache map[string]Wc.Dictionary) http.HandlerFunc {

	// for now I'll just assume it's a POST - assuming it's within 10MB, too
	return func(rw http.ResponseWriter, request *http.Request) {

		switch request.Method {
		case "GET":
			func() {
				if filenames, err := FilenamesInRequest(request); err != nil {

					// eg, http://localhost:3000/files - return the cached filenames
					known_files := make([]string, 0, len(WcCache))
					for filename := range WcCache {
						known_files = append(known_files, filename)
					}
					emit_json(rw, known_files)

				} else {

					if len(filenames) == 0 {
						// eg, http://localhost:3000/?filename=foo - return WcCache["foo"]
						emit_json(rw, WcCache[filenames[0]])
					} else {
						// eg, http://localhost:3000/?filename=foo&filename=bar
						OutCache := Wc.Dictionary{Total: 0, Words: map[string]int{}}

						for _, filename := range filenames {
							OutCache.Total += WcCache[filename].Total

							for k, v := range WcCache[filename].Words {
								if _, ok := OutCache.Words[k]; ok {
									OutCache.Words[k] += v
								} else {
									OutCache.Words[k] = v
								}
							}
						}

						emit_json(rw, OutCache)
					}
				}

			}()

		case "DELETE":
			func() {

				filenames, _ := FilenamesInRequest(request)
				fmt.Printf("got a DELETE %s\n", filenames[0])
				delete(WcCache, filenames[0])
				emit_json(rw, []byte(nil))

			}()

		case "POST":
			func() {
				file, _, err := request.FormFile("file")

				filenames, _ := FilenamesInRequest(request)

				// bail out if it's in WcCache already
				if seen, ok := WcCache[filenames[0]]; ok {
					fmt.Printf("Returning cached data for %s\n", filenames[0])
					emit_json(rw, seen)
					return
				}

				if err == nil {
					// slurp the file, then convert byte array to a string
					corpus_bytes, _ := ioutil.ReadAll(file)
					corpus := string(corpus_bytes[:])

					// split into an array of words, then count them
					corpus_words := Wc.SplitWords(corpus)
					pages, _ := WcPaginate.MakePageSlices(len(corpus_words), 10)

					countsC := make(chan Wc.Dictionary)
					for low, high := range pages {
						low, high := low, high
						go func() {
							countsC <- Wc.CountWords(corpus_words[low:high])
						}()
					}

					var gotPages []Wc.Dictionary
					for i := 0; i < len(pages); i++ {
						gotPages = append(gotPages, <-countsC)
					}

					counts := sum_counts(gotPages)

					// store in WcCache
					fmt.Printf("storing data for %s in WcCache\n", filenames[0])
					WcCache[filenames[0]] = counts

					// send response
					emit_json(rw, counts)
				} else {
					fmt.Fprintf(rw, "encountered error: %s", err)
				}
			}()
		}
	}
}

func sum_counts(dicts []Wc.Dictionary) Wc.Dictionary {
	OutCache := Wc.Dictionary{Total: 0, Words: map[string]int{}}

	for _, d := range dicts {
		OutCache.Total += d.Total

		for k, v := range d.Words {
			if _, ok := OutCache.Words[k]; ok {
				OutCache.Words[k] += v
			} else {
				OutCache.Words[k] = v
			}
		}
	}

	return OutCache
}

func HandleAdminFiles(WcCache map[string]Wc.Dictionary) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		switch request.Method {

		case "DELETE":
			func() {
				for k := range WcCache {
					delete(WcCache, k)
				}
				emit_json(rw, []byte(`[]`))
			}()

		case "GET":
			func() {
				emit_json(rw, WcCache)
			}()
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

	// get a cache
	WcCache := InitCache()

	// register handlers and start listening for requests
	http.HandleFunc("/files", HandleWcFile(WcCache))
	http.HandleFunc("/admin/files", HandleAdminFiles(WcCache))
	http.ListenAndServe(portString, nil)
}
