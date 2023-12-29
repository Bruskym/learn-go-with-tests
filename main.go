package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type playersStore interface {
	getScore(name string) int
	storeWin(name string)
	getLeague() []Player
}

type Player struct {
	Name string
	Wins int
}

type playerServer struct {
	Store  playersStore
	Router http.Handler
}

func newPlayerServer(store playersStore) *playerServer {
	router := http.NewServeMux()
	p := new(playerServer)
	p.Store = store

	router.Handle("/players/", http.HandlerFunc(p.playersHandle))
	router.Handle("/league/", http.HandlerFunc(p.leagueHandle))

	p.Router = router

	return p
}

func (p *playerServer) playersHandle(w http.ResponseWriter, r *http.Request) {
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

func (p *playerServer) leagueHandle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	enconder := json.NewEncoder(w)

	if err := enconder.Encode(p.getLeagueTable()); err != nil {
		log.Fatal(err)
	}
}

func (p *playerServer) getLeagueTable() []Player {
	return p.Store.getLeague()
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

func (i *InMemoryStorePlayers) getLeague() []Player {
	var league []Player

	for player, wins := range i.score {
		league = append(league, Player{player, wins})
	}

	return league
}

func main() {
	server := newPlayerServer(NewInMemoryStorePlayers())

	if err := http.ListenAndServe(":8080", server.Router); err != nil {
		log.Fatal(err)
	}
}
