package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/callback", callbackHandler)

	http.ListenAndServe("localhost:7002", nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// gather information about request and log it
	uri := r.URL.String()
	method := r.Method
	fmt.Printf("uri = %s, method = %s", uri, method)
	w.WriteHeader(http.StatusNotImplemented)
}
