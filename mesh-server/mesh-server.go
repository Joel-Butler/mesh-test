package main

//web server derived from https://golang.org/doc/articles/wiki/

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func main() {
	//define environment variable
	http.HandleFunc("/", handler)
	meshServiceUrl = os.Getenv("JB_MESH_SERVICE")

	if meshServiceUrl == "" {
		meshServiceUrl = "http://localhost:8081"
	}

	listenPort = os.Getenv("JB_MESH_SERVER_PORT")
	if listenPort == "" {
		listenPort = ":8080"
	}

	log.Fatal(http.ListenAndServe(listenPort, nil))
}
