package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

type InMemoryPlayerStore struct {
	sync.Mutex
	store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		store: make(map[string]int),
	}
}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) int {
	s.Lock()
	defer s.Unlock()
	return s.store[name]
}

func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.Lock()
	defer s.Unlock()
	old, ok := s.store[name]
	if !ok {
		s.store[name] = 1
		return
	}
	s.store[name] = old + 1
}

func (s *InMemoryPlayerStore) GetLeague() []Player {
	players := make([]Player, 0)
	for k, v := range s.store {
		players = append(players, Player{Name: k, Wins: v})
	}
	return players
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))

	router.Handle("/players/", http.HandlerFunc(p.playerHandler))
	p.Handler = router
	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.getLeagueTable())

}

func (p *PlayerServer) getLeagueTable() []Player {
	return p.store.GetLeague()
}

func (p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]
	switch r.Method {
	case http.MethodGet:
		p.showScore(w, player)
	case http.MethodPost:
		p.processWin(w, player)
	}
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {

	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}
