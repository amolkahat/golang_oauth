package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amolkahat/golang_oauth/controllers"
	"github.com/joho/godotenv"
)

func main() {

	// Create simple server and run

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	server := http.Server{
		Addr:    fmt.Sprintf(":8000"),
		Handler: controllers.New(),
	}

	log.Printf("Starting server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}
