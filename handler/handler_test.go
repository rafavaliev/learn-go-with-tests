package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
	old, ok := s.scores[name]
	if !ok {
		s.scores[name] = 1
		return
	}
	s.scores[name] = old + 1

}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		winCalls: make([]string, 0),
	}
	server := &PlayerServer{store: store}
	table := []struct {
		title            string
		playerName       string
		want             string
		expectStatusCode int
	}{
		{
			"returns Pepper's score",
			"Pepper",
			"20",
			200,
		},
		{
			"returns Floyd's score",
			"Floyd",
			"10",
			200,
		},
		{
			"return 404 on missing players",
			"Appolo",
			"0",
			404,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/players/"+tt.playerName, nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseBody(t, response.Body.String(), tt.want)
			assertResponseCode(t, response.Code, tt.expectStatusCode)

		})
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		make([]string, 0),
	}
	server := &PlayerServer{&store}

	table := []struct {
		title            string
		playerName       string
		want             string
		expectStatusCode int
	}{
		{
			"return status accepted on post requests",
			"Appolo",
			"1",
			http.StatusAccepted,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/players/"+tt.playerName, nil)
			response := httptest.NewRecorder()

			checkNotFound(t, server, tt.playerName)

			request, _ = http.NewRequest(http.MethodPost, "/players/"+tt.playerName, nil)
			response = httptest.NewRecorder()

			server.ServeHTTP(response, request)
			assertResponseCode(t, response.Code, tt.expectStatusCode)

			checkFoundWithBody(t, server, tt.playerName, tt.want)

			if len(store.winCalls) != 1 {
				t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
			}
		})
	}
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertResponseCode(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}

func newPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func checkNotFound(t *testing.T, server *PlayerServer, name string) {
	t.Helper()
	request, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	assertResponseCode(t, response.Code, http.StatusNotFound)
}

func checkFoundWithBody(t *testing.T, server *PlayerServer, name, want string) {
	t.Helper()
	request, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)
	assertResponseCode(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), want)
}

func assertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
