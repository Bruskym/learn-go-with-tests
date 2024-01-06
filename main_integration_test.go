package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterWinAndListPlayer(t *testing.T) {
	player := "Antonio"

	database, clearDatabase := createTempFile(t, "")
	defer clearDatabase()

	store := NewFileSystemStore(database)
	server := NewPlayerServer(store)

	server.Router.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))
	server.Router.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))
	server.Router.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))

	t.Run("get Score", func(t *testing.T) {
		req := newGetScoreRequest(player)
		res := httptest.NewRecorder()
		server.Router.ServeHTTP(res, req)

		got := res.Body.String()
		want := "3"

		assertStatusCode(t, res.Code, http.StatusOK)
		assertBodyContent(t, want, got)
	})

	t.Run("get League", func(t *testing.T) {
		req := newGetLeagueRequest()
		res := httptest.NewRecorder()

		server.Router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusOK)

		want := []Player{
			{"Antonio", 3},
		}
		got := getLeagueFromBody(t, res.Body)

		assertResponseContentType(t, res.Header().Get("content-type"), "application/json")
		assertLeague(t, want, got)

	})

}
