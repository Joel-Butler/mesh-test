package main

//web server derived from https://golang.org/doc/articles/wiki/

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var meshServiceUrl string
var listenPort string

func handler(w http.ResponseWriter, r *http.Request) {
	//poll response from back-end service
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(meshServiceUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Hi there, your random number is %s!", body)
	log.Printf("Request: %s client: %s forwarded for: %s", r.URL.Path, r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go Health status server")
	w.WriteHeader(200)
}

func main() {
	//define environment variable
	http.HandleFunc("/", handler)
	meshServiceUrl = os.Getenv("JB_MESH_SERVICE")

	if meshServiceUrl == "" {
		meshServiceUrl = "http://localhost:8081/api"
	}

	wg := new(sync.WaitGroup)

	wg.Add(2)

	healthServer := http.NewServeMux()
	healthServer.HandleFunc("/health", healthHandler)

	apiserver := http.NewServeMux()
	apiserver.HandleFunc("/", handler)

	listenPort = os.Getenv("JB_MESH_SERVER_PORT")
	if listenPort == "" {
		listenPort = ":8080"
	}

	go func() {
		server2 := http.Server{
			Addr:    listenPort,
			Handler: apiserver,
		}
		log.Fatal(server2.ListenAndServe())
		wg.Done()
	}()

	go func() {
		server1 := http.Server{
			Addr:    ":8079", // :{port}
			Handler: healthServer,
		}
		server1.ListenAndServe()
		wg.Done()
	}()
	wg.Wait()
}
