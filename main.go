package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// configure the songs directory name and port

	const videoDir = "videos"

	const port = 8080

	// add logging everytime a video request is received and content is served

	http.Handle("/", addHeaders(http.FileServer(http.Dir(videoDir))))

	fmt.Printf("Starting server on %v\n", port)

	log.Printf("Serving %s on HTTP port: %v\n", videoDir, port)

	// serve and log errors

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// addHeaders will act as middleware to give us CORS support

func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		h.ServeHTTP(w, r)
	}
}
