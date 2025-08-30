package main

import (
	"net/http"
	"log"
)

func main() {
	log.Println("Server started")
	http.ListenAndServe(":8081", nil)
}