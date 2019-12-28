// Package main illustrates how to handle errors from goroutines
package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

// make a request to many urls with different goroutines
// check the response from making a request to the slice of strings
func checkStatus(done <-chan interface{}, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)

		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)
			// create the response
			result = Result{Error: err, Response: resp}

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
	errCount := 0
	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://lumexralph.github.io", "https://local", "https://lo", "https://loda"}
	for result := range checkStatus(done, urls...) {
		// the main goroutine can better handle the error
		// it is the goroutine that spawned the other goroutine
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, ending the process")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}

}
