package main

import (
	"net/http"
	"fmt"
	"time"
)

const urlRedirector = "https://fresh-argon-152122.appspot.com/redirect"

func main() {
	for {
		go callRedirect()
		time.Sleep(1000 * time.Millisecond)
	}
}

func callRedirect() {
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
		fmt.Printf("Status code: %d, Redirect location: %s\n", resp.StatusCode, resp.Header.Get("Location"))

		if resp.StatusCode < 300 || resp.StatusCode >= 400 {
			break
		}

		urlString = resp.Header.Get("Location")
	}
}
