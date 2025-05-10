package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", "8080", "Port to serve on")
	dir := flag.String("dir", ".", "Directory to serve")
	flag.Parse()

	// Check if directory exists
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		fmt.Printf("hserve: directory '%s' does not exist\n", *dir)
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	fmt.Printf("hserve: serving %s on http://localhost:%s\n", *dir, *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		fmt.Printf("hserve: error: %v\n", err)
		os.Exit(1)
	}
}
