package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("Server starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("err")
	}
}
