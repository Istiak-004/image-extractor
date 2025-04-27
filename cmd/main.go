package main

import (
	"log"
	"os"

	"github.com/istiak-004/image-extractor/internals/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	server := app.NewServer()
	log.Printf("Server starting on port %s", port)
	log.Fatal(server.Start(":" + port))
}
