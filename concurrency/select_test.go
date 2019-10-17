package concurrency

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {

	t.Run("Get fastest url", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)
		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, _ := Racer(fastURL, slowURL)

		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}

	})

	t.Run("Error, if servers didn't respond withint timeout", func(t *testing.T) {
		timeout := 10 * time.Millisecond
		server := makeDelayedServer(timeout * 2)
		defer server.Close()

		_, err := TimeoutableRacer(server.URL, server.URL, timeout)

		if err == nil {
			t.Error("Expected error, but got nothing")
		}
	})
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
