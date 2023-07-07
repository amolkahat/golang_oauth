package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amolkahat/golang_oauth/base"
)

func main() {

	// Create simple server and run

	server := http.Server{
		Addr:    fmt.Sprintf(":8000"),
		Handler: base.New(),
	}

	log.Printf("Starting server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}
