package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func (l League) FindPlayerByName(name string) *Player {
	for i, player := range l {
		if player.Name == name {
			return &l[i]
		}
	}
	return nil
}

func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player

	decode := json.NewDecoder(rdr)

	if err := decode.Decode(&league); err != nil {
		return nil, fmt.Errorf("error parsing JSON content from the database: %w", err)
	}

	return league, nil
}
