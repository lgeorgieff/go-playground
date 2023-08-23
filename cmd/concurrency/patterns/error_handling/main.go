package main

import (
	"fmt"
	"net/http"
	"os"
)

type Result struct {
	Response *http.Response
	Error    error
}

func CheckStatus(done <-chan struct{}, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)

		for _, url := range urls {
			resp, err := http.Get(url)
			result := Result{
				Response: resp,
				Error:    err,
			}
			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()

	return results
}

func main() {
	const maxErrorThreshold = 3

	done := make(chan struct{})
	defer close(done)

	urls := []string{
		"https://google.com",
		"https://doesnotexist-1",
		"https://spiegel.de",
		"https://doesnotexist-2",
		//"https://doesnotexist-3",
	}
	var errCounter uint
	for result := range CheckStatus(done, urls...) {
		if result.Error != nil {
			fmt.Fprintln(os.Stderr, result.Error)

			errCounter++

			if errCounter >= maxErrorThreshold {
				fmt.Fprintf(os.Stderr, "too many errors (%d/%d)\n", errCounter, maxErrorThreshold)
				break
			}

			continue
		}
		fmt.Printf("Response for %s resulted in %v\n", result.Response.Request.URL, result.Response.Status)
	}
}
