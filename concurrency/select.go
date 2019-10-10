package concurrency

import (
	"fmt"
	"net/http"
	"time"
)

var tenSecondTimeout = 10 * time.Second

func Racer(aUrl, bUrl string) (winner string, error error) {
	return TimeoutableRacer(aUrl, bUrl, tenSecondTimeout)
}

func TimeoutableRacer(aUrl, bUrl string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(aUrl):
		return aUrl, nil
	case <-ping(bUrl):
		return bUrl, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout while waiting for %s and %s", aUrl, bUrl)
	}
}

func ping(url string) chan bool {
	ch := make(chan bool)
	go func() {
		http.Get(url)
		ch <- true
	}()
	return ch
}
