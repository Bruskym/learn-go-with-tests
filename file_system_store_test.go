package poker

import (
	"os"
	"testing"
)

func TestFileStore(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name":"Antonio", "Wins":6},
		{"Name":"Ellen", "Wins":2}
	]`)

	defer cleanDatabase()
	store, err := NewFileSystemStore(database)
	assertError(t, err)

	t.Run("create store with empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemStore(database)

		assertError(t, err)
	})

	t.Run("get league from database", func(t *testing.T) {

		got := store.getLeague()
		want := []Player{
			{"Antonio", 6},
			{"Ellen", 2},
		}

		assertLeague(t, want, got)

		got = store.getLeague()
		assertLeague(t, want, got)

	})

	t.Run("get player from database", func(t *testing.T) {
		got := store.getScore("Antonio")
		want := 6

		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing player", func(t *testing.T) {
		store.storeWin("Antonio")

		got := store.getScore("Antonio")
		want := 7

		assertScoreEquals(t, want, got)
	})

	t.Run("store wins for new player", func(t *testing.T) {
		store.storeWin("Joao")

		got := store.getScore("Joao")
		want := 1

		assertScoreEquals(t, want, got)

		store.storeWin("Joao")
	})
}

func assertScoreEquals(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("The expected score was %d and was returned %d", want, got)
	}
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file. %v", err)
	}

	tempFile.Write([]byte(initialData))

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile

}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("No errors were expected, but %v", err)
	}
}
