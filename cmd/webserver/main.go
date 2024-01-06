package main

import (
	poker "API-Rest/Go-with-Tests"
	"log"
	"net/http"
	"os"
)

const storeFileName = "db.json"

func main() {
	file, err := os.OpenFile(storeFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("unable to open the file. %v", err)
	}

	store := poker.NewFileSystemStore(file)
	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":8000", server.Router); err != nil {
		log.Fatal(err)
	}
}
