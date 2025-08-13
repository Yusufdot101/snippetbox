package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// create a file serever which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// use the Handle() method to register the file server as the handler for
	// all url paths that start with "/static'". for matching paths, we stripp the
	// "/static" prefix before the request reacehs teh file server
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// register other routes
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snippet/view", ViewSnippet)
	mux.HandleFunc("/snippet/create", CreateSnippet)

	port := ":3000"
	fmt.Printf("Server listening on port: %v", port)
	http.ListenAndServe(port, mux)
}
