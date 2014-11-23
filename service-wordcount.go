package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	// "regexp"
)

func wc_file(rw http.ResponseWriter, request *http.Request) {
	// request_json, _ := json.Marshal(request)
	// rw.Write(request_json)

	file, _, err := request.FormFile("file")
	if err == nil {
		content, _ := ioutil.ReadAll(file)
		fmt.Printf("got data: %s", content)
	} else {
		fmt.Fprintf(rw, "encoutered error: %s", err)
	}

	fmt.Fprintln(rw, "Success!")
}

func main() {
	http.HandleFunc("/", wc_file)
	http.ListenAndServe(":3000", nil)
}
