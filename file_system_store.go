package poker

import (
	"encoding/json"
	"os"
)

type FileSystemStore struct {
	data   json.Encoder
	league League
}

func NewFileSystemStore(database *os.File) *FileSystemStore {
	database.Seek(0, 0)

	league, _ := NewLeague(database)
	return &FileSystemStore{*json.NewEncoder(&Tape{database}), league}
}

func (f *FileSystemStore) getLeague() League {
	return f.league
}

func (f *FileSystemStore) getScore(name string) int {

	league := f.getLeague()

	if player := league.FindPlayerByName(name); player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemStore) storeWin(name string) {

	if player := f.league.FindPlayerByName(name); player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.data.Encode(f.league)
}
