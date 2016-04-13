package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received HTTP %s request to %s for %v from %s \n", r.Method, r.Host, r.URL, r.RemoteAddr)
	fmt.Fprintf(w, "Thank you and good bye.\n")
}

func main() {
	portEnv := os.Getenv("PORT")
	port, err := strconv.ParseInt(portEnv, 10, 0)
	if err != nil {
		log.Fatalf("Failed to parse PORT: %v", err)
	}
	fmt.Printf("Server registering to any interface on port %d starting..\n", port)
	http.HandleFunc("/", handle)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
