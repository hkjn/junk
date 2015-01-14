// app.go: A trivial Go service to deploy on CoreOS and keep up-to-date with CD.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Printf("app.go starting in stage %s..\n", os.Getenv("STAGE"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi from app.go!\n")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
