package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("port", 8081, "port to start demo server")

type DemoServer struct{}

func (d *DemoServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println("Hit me " + fmt.Sprint(*port))
	rw.WriteHeader(200)
	rw.Write([]byte("All good from server " + fmt.Sprint(*port) + "\n"))
}

func main() {
	flag.Parse()

	demo := &DemoServer{}
	log.Fatalln(http.ListenAndServe(":"+fmt.Sprint(*port), demo))
}
