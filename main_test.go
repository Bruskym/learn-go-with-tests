package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type stubPlayerScore struct {
	score           map[string]int
	registerWinCall []string
	league          []Player
}

func (s *stubPlayerScore) getScore(name string) int {
	return s.score[name]
}

func (s *stubPlayerScore) storeWin(name string) {
	s.registerWinCall = append(s.registerWinCall, name)
}

func (s *stubPlayerScore) getLeague() []Player {
	return s.league
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
	wantedLeague := []Player{
		{"Antonio", 5},
		{"Ellen", 3},
	}

	testScore := &stubPlayerScore{league: wantedLeague}
	server := newPlayerServer(testScore)

	t.Run("It returns 200 status code on /league/", func(t *testing.T) {
		req := newGetLeagueRequest()
		res := httptest.NewRecorder()

		server.Router.ServeHTTP(res, req)

		var got []Player

		decoder := json.NewDecoder(res.Body)

		if err := decoder.Decode(&got); err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", res.Body, err)
		}

		assertStatusCode(t, res.Code, http.StatusOK)
	})

	t.Run("It returns the league table as JSON", func(t *testing.T) {
		req := newGetLeagueRequest()
		res := httptest.NewRecorder()

		server.Router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusOK)

		got := getLeagueFromBody(t, res.Body)

		assertLeague(t, wantedLeague, got)

		want := "application/json"

		assertResponseContentType(t, res.Result().Header.Get("content-type"), want)

	})
}

func assertLeague(t *testing.T, want, got []Player) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want : %v. Got : %v", want, got)
	}
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

func assertResponseContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response did not have content-type of %s, got %v", want, got)
	}
}

func newGetScoreRequest(name string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetLeagueRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/league/", nil)
	return req
}

func getLeagueFromBody(t *testing.T, body io.Reader) []Player {
	var league []Player

	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&league); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return league
}

func newPostScoreRequest(name string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}
