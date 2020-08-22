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
	"sync"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", rand.Intn(100000))
	log.Printf("Request: %s client: %s forwarded for: %s", r.URL.Path, r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"))
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go Health status server")
	w.WriteHeader(200)
}

func main() {
	http.HandleFunc("/", handler)
	rand.Seed(time.Now().UTC().UnixNano())
	//lazy option for a second listener for health checks - from https://medium.com/rungo/running-multiple-http-servers-in-go-d15300f4e59f

	wg := new(sync.WaitGroup)

	wg.Add(2)

	healthServer := http.NewServeMux()
	healthServer.HandleFunc("/health", healthHandler)
	apiserver := http.NewServeMux()
	apiserver.HandleFunc("/api", handler)

	go func() {
		server2 := http.Server{
			Addr:    ":8081", // :{port}
			Handler: apiserver,
		}
		log.Fatal(server2.ListenAndServe())
		wg.Done()
	}()

	go func() {
		server1 := http.Server{
			Addr:    ":8082", // :{port}
			Handler: healthServer,
		}
		server1.ListenAndServe()
		wg.Done()
	}()

	wg.Wait()
}
