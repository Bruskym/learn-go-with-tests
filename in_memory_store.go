package main

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
