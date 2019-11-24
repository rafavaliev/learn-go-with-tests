package handler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
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

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
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
