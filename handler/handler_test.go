package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const jsonContentType = "application/json"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
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

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		winCalls: make([]string, 0),
	}
	server := NewPlayerServer(store)
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
	store := &StubPlayerStore{
		map[string]int{},
		make([]string, 0),
		nil,
	}
	server := NewPlayerServer(store)

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
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertResponseCode(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}

func TestLeague(t *testing.T) {
	store := &StubPlayerStore{}
	server := NewPlayerServer(store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		_ = getLeagueFromResponse(t, response.Body)
		assertResponseCode(t, response.Code, http.StatusOK)
	})

	t.Run("it returns the league table as Json", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store = &StubPlayerStore{nil, nil, wantedLeague}
		server = NewPlayerServer(store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)

	})
}

func TestRecordingWinsAndRetrievingThemWithInMemoryStore(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertResponseCode(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertResponseCode(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %v, got %v", want, response.Result().Header.Get("content-type"))
	}
}

func assertLeague(t *testing.T, got []Player, wantedLeague []Player) {

	t.Helper()
	if !reflect.DeepEqual(wantedLeague, got) {
		t.Errorf("Got %v, want %v", got, wantedLeague)
	}
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Player {
	t.Helper()
	var got []Player
	err := json.NewDecoder(body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response for server %q into slice of Players, %v", body, err)
	}
	return got
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
