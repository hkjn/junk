package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var maxMb int64 = 1

func receive(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received HTTP %s on /\n", r.Method)
	if r.Method != "PUT" {
		fmt.Fprintf(w, "Bad method: %s", r.Method)
		return
	}
	f, err := os.Create("outputfile")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()
	b := make([]byte, maxMb*1024*1024)
	total := 0
	for {
		done := false
		n, err := r.Body.Read(b)
		if err != nil {
			if err == io.EOF {
				done = true
			} else {
				log.Fatalf("failed to read: %v", err)
			}
		}
		total += n
		fmt.Printf("read %d bytes (%d total, EOF=%v)\n", n, total, done)
		n, err = f.Write(b[:n])
		if err != nil {
			log.Fatalf("failed to write: %v", err)
		}
		fmt.Printf("wrote %d bytes\n", n)
		if done {
			break
		}
	}
	fmt.Printf("Received %d bytes from client.\n", total)
	fmt.Fprintf(w, "Received %d bytes. Thanks!\n", total)
}

func main() {
	maxMbEnv := os.Getenv("MAX_BYTES_MB")
	var err error
	maxMb, err = strconv.ParseInt(maxMbEnv, 10, 0)
	if err != nil {
		log.Fatalf("Failed to parse MAX_BYTES_MB: %v", err)
	}
	fmt.Printf("Server accepting max HTTP PUT of %d MB starting..\n", maxMb)
	http.HandleFunc("/", receive)
	if err := http.ListenAndServe(":443", nil); err != nil {
		panic(err)
	}
}
