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
	for {
		go execRead()
		time.Sleep(100 * time.Millisecond)
	}
}

func execRead() {
	resp, err := http.Get(sensitive.AppEngineUrl)
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
