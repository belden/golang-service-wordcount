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
		rw.Write([]byte(`{"error": "error in json encoding"}`))
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

func InitCache() map[string]Wc.Dictionary {
	return make(map[string]Wc.Dictionary)
}

func HandleWcFile(WcCache map[string]Wc.Dictionary) http.HandlerFunc {

	// for now I'll just assume it's a POST - assuming it's within 10MB, too
	return func(rw http.ResponseWriter, request *http.Request) {

		switch request.Method {
		case "GET":
			func() {
				if filename, err := FilenameInRequest(request); err != nil {

					// eg, http://localhost:3000/files - return the cached filenames
					filenames := make([]string, 0, len(WcCache))
					for filename := range WcCache {
						filenames = append(filenames, filename)
					}
					emit_json(rw, filenames)

				} else {

					// eg, http://localhost:3000/?filename=foo - return WcCache["foo"]
					emit_json(rw, WcCache[filename])

				}

			}()

		case "DELETE":
			func() {

				filename, _ := FilenameInRequest(request)
				fmt.Printf("got a DELETE %s\n", filename)
				delete(WcCache, filename)
				emit_json(rw, []byte(nil))

			}()

		case "POST":
			func() {
				file, _, err := request.FormFile("file")

				filename, _ := FilenameInRequest(request)

				// bail out if it's in WcCache already
				if seen, ok := WcCache[filename]; ok {
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

					// store in WcCache
					fmt.Printf("storing data for %s in WcCache\n", filename)
					WcCache[filename] = counts

					// send response
					emit_json(rw, counts)
				} else {
					fmt.Fprintf(rw, "encountered error: %s", err)
				}
			}()
		}
	}
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
