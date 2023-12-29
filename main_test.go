package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubPlayerScore struct {
	score           map[string]int
	registerWinCall []string
}

func (s *stubPlayerScore) getScore(name string) int {
	return s.score[name]
}

func (s *stubPlayerScore) storeWin(name string) {
	s.registerWinCall = append(s.registerWinCall, name)
}

func TestGetScore(t *testing.T) {

	testStore := &stubPlayerScore{
		score: map[string]int{
			"Antonio": 20,
			"Jack":    10,
		},
	}

	testServer := newPlayerServer(testStore)
	t.Run("Antonio's score is returned correctly", func(t *testing.T) {
		req := newGetScoreRequest("Antonio")
		res := httptest.NewRecorder()

		testServer.Router.ServeHTTP(res, req)

		got := res.Body.String()
		want := "20"

		assertStatusCode(t, http.StatusOK, res.Code)
		assertBodyContent(t, got, want)
	})
	t.Run("Jack's score is returned correctly", func(t *testing.T) {
		req := newGetScoreRequest("Jack")
		res := httptest.NewRecorder()

		testServer.Router.ServeHTTP(res, req)

		got := res.Body.String()
		want := "10"

		assertStatusCode(t, http.StatusOK, res.Code)
		assertBodyContent(t, got, want)
	})

	t.Run("It returns a 404 status code if the player doesn't exist", func(t *testing.T) {
		req := newGetScoreRequest("Test")
		res := httptest.NewRecorder()

		testServer.Router.ServeHTTP(res, req)

		got := res.Code
		want := http.StatusNotFound

		assertStatusCode(t, want, got)
	})
}

func TestPostWin(t *testing.T) {
	player := "Antonio"
	t.Run("It returns accept status code when POST", func(t *testing.T) {
		server := newPlayerServer(&stubPlayerScore{})

		req := newPostScoreRequest(player)
		res := httptest.NewRecorder()

		server.Router.ServeHTTP(res, req)

		got := res.Code
		want := http.StatusAccepted

		assertStatusCode(t, want, got)
	})

	t.Run("if when making a POST the store is called correctly", func(t *testing.T) {
		store := &stubPlayerScore{}
		server := newPlayerServer(store)

		req := newPostScoreRequest(player)
		res := httptest.NewRecorder()

		server.Router.ServeHTTP(res, req)

		got := len(store.registerWinCall)
		want := 1

		if got != want {
			t.Errorf("the function was expected to be called %d times, but it was called %d times", want, got)
		}

		if store.registerWinCall[0] != player {
			t.Errorf("who should have been registered as the winner was %q, but instead it was %q", player, store.registerWinCall[0])
		}

	})
}

func TestLeagueRoute(t *testing.T) {
	testScore := &stubPlayerScore{}
	server := newPlayerServer(testScore)

	t.Run("It returns 200 status code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/league/", nil)
		res := httptest.NewRecorder()

		server.Router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusOK)
	})
}
func assertBodyContent(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("The expected score was %q and was returned %q", want, got)
	}
}

func assertStatusCode(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("Returned status code %d was expected status code %d", got, want)
	}
}

func newGetScoreRequest(name string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostScoreRequest(name string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}
