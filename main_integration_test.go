package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterWinAndListPlayer(t *testing.T) {
	player := "Antonio"
	store := NewInMemoryStorePlayers()

	server := newPlayerServer(store)

	server.Router.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))
	server.Router.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))

	// GET
	req := newGetScoreRequest(player)
	res := httptest.NewRecorder()
	server.Router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "2"

	assertStatusCode(t, res.Code, http.StatusOK)
	assertBodyContent(t, want, got)
}
