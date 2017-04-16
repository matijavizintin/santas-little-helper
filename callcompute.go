package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"santas-little-helper/sensitive"
	"time"
)

func main() {
	throttleChan := make(chan int, 10)

	for {
		throttleChan <- 1
		go execCall(throttleChan)
		time.Sleep(10 * time.Millisecond)
	}
}

func execCall(throttleChan chan int) {
	defer func() {
		<-throttleChan
	}()

	resp, err := http.Get(sensitive.LoadBalancerUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Stdout.WriteString(fmt.Sprintf("Status code: %v Content: ", resp.StatusCode))
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	os.Stdout.Write([]byte{'\n'})
}
