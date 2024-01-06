package poker

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileSystemStore struct {
	data   json.Encoder
	league League
}

func NewFileSystemStore(database *os.File) (*FileSystemStore, error) {
	err := initialiseDBFile(database)

	if err != nil {
		return nil, fmt.Errorf("error getting information about file %v", err)
	}

	league, err := NewLeague(database)

	if err != nil {
		return nil, fmt.Errorf("problem loading file player storage: %v", err)
	}

	return &FileSystemStore{*json.NewEncoder(&Tape{database}), league}, nil
}

func initialiseDBFile(dbFile *os.File) error {
	dbFile.Seek(0, 0)

	fileStats, err := os.Stat(dbFile.Name())

	if err != nil {
		return fmt.Errorf("error getting information about file %v", err)
	}

	if fileStats.Size() == 0 {
		dbFile.Write([]byte("[]"))
		dbFile.Seek(0, 0)
	}

	return nil
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
