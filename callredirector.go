package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"santas-little-helper/sensitive"
	"time"
)

const printOnlyErrors = true

func main() {
	throttleChan := make(chan int, 10)

	c := http.DefaultClient
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for {
		throttleChan <- 1
		go callRedirect(c, throttleChan)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
}

func callRedirect(c *http.Client, throttleChan chan int) {
	defer func() {
		<-throttleChan
	}()

	urlString := sensitive.AppEngineUrl
	for {
		r, err := http.NewRequest("GET", urlString, nil)
		if err != nil {
			fmt.Println("Error creating get request:", err)
			return
		}
		resp, err := c.Do(r)
		if err != nil {
			fmt.Println("Error doing request: ", err)
			return
		}
		resp.Body.Close()

		if !printOnlyErrors || resp.StatusCode >= 400 {
			fmt.Printf("Status code: %d, Response from: %s Redirect location: %s\n", resp.StatusCode, urlString, resp.Header.Get("Location"))
		}

		if resp.StatusCode < 300 || resp.StatusCode >= 400 {
			break
		}

		urlString = resp.Header.Get("Location")
	}
}
