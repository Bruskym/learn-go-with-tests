package main

import (
	"log"
	"net/http"
)

func main() {
	server := newPlayerServer(NewInMemoryStorePlayers())

	if err := http.ListenAndServe(":8080", server.Router); err != nil {
		log.Fatal(err)
	}
}
