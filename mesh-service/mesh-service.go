package main

/*web server derived from https://golang.org/doc/articles/wiki/
 *
 * this service simply returns a random number to any web query on port 8080
 */

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", rand.Intn(100000))
	log.Printf("Request: %s client: %s forwarded for: %s", r.URL.Path, r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"))
}

func main() {
	http.HandleFunc("/", handler)
	rand.Seed(time.Now().UTC().UnixNano())
	log.Fatal(http.ListenAndServe(":8081", nil))
}
