package concurrency

import (
	"net/http"
	"time"
)

func Racer(aUrl, bUrl string) string {
	aDuration := measureResponseTime(aUrl)
	bDuration := measureResponseTime(bUrl)

	if aDuration < bDuration {
		return aUrl
	}
	return bUrl

}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
