package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8888", "server address:port")
	flag.Parse()
	srv := http.FileServer(http.Dir("."))
	log.Printf("listening on %q...", *addr)
	log.Fatal(http.ListenAndServe(*addr, srv))
}
