package main

import (
	"net/http"
	"fmt"
	"io"
	"os"
	"time"
)

const urlRandom = "https://fresh-argon-152122.appspot.com/random"

func main() {
	for {
		go execRead()
		time.Sleep(100 * time.Millisecond)
	}
}

func execRead() {
	resp, err := http.Get(urlRandom)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		os.Stdout.WriteString(fmt.Sprintf("Status code: %v Content: ", resp.StatusCode))
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	os.Stdout.Write([]byte{'\n'})
}
