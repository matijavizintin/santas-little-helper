package main

import (
	"net/http"
	"fmt"
	"time"
	"math/rand"
)

const urlRedirector = "https://fresh-argon-152122.appspot.com/redirect"

func main() {
	throttleChan := make(chan int, 50)

	for {
		throttleChan <- 1
		go callRedirect(throttleChan)
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	}
}

func callRedirect(throttleChan chan int) {
	defer func() {
		<- throttleChan
	}()

	c := http.DefaultClient
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	urlString := urlRedirector
	for {
		r, err := http.NewRequest("GET", urlString, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := c.Do(r)
		if err != nil {
			fmt.Println(err)
			return
		}
		if r.Body != nil {
			r.Body.Close()
		}

		if resp.StatusCode >= 400 {
			fmt.Printf("Status code: %d, Response from: %s Redirect location: %s\n", resp.StatusCode, urlString, resp.Header.Get("Location"))
		}

		if resp.Body != nil {
			resp.Body.Close()
		}

		if resp.StatusCode < 300 || resp.StatusCode >= 400 {
			break
		}

		urlString = resp.Header.Get("Location")
	}
}
