package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type playersStore interface {
	getScore(name string) int
	storeWin(name string)
}

type playerServer struct {
	Store playersStore
}

func (p *playerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodGet:
		p.showScore(w, player)
	case http.MethodPost:
		p.registerWin(w, player)
	}
}

func (p *playerServer) showScore(w http.ResponseWriter, player string) {
	score := p.Store.getScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *playerServer) registerWin(w http.ResponseWriter, player string) {
	p.Store.storeWin(player)
	w.WriteHeader(http.StatusAccepted)
}

type InMemoryStorePlayers struct {
	score map[string]int
}

func NewInMemoryStorePlayers() *InMemoryStorePlayers {
	return &InMemoryStorePlayers{map[string]int{}}
}

func (i *InMemoryStorePlayers) getScore(name string) int {
	return i.score[name]
}

func (i *InMemoryStorePlayers) storeWin(name string) {
	i.score[name]++
}

func main() {
	server := &playerServer{NewInMemoryStorePlayers()}

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
