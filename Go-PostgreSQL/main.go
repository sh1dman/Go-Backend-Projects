package main

import (
	"Go-PostgreSQL/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Starting Server on the port 8080...")

	log.Fatal(http.ListenAndServe(8080, r))
}
