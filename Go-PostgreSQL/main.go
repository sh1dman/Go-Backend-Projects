package main

import (
	"Go-PostgreSQL/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	// Use a different port (e.g., 8081)
	port := ":8081"
	fmt.Printf("Starting server on the port %s...\n", port[1:])

	log.Fatal(http.ListenAndServe(port, r))
}
