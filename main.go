package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	sigs := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	msg := make(chan string, 1)
	go func() {
		// Receive input in a loop
		for {
			var s string
			fmt.Scan(&s)
			// Send what we read over the channel
			msg <- s
		}
	}()

loop:
	for {
		select {
		case <-sigs:
			fmt.Println("Got shutdown, exiting")
			// Break out of the outer for statement and end the program
			break loop
		case s := <-msg:
			time_start := time.Now()
			resp, err := http.Get(s)
			if err != nil {
				fmt.Printf("error")
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("error")
			}
			l := len(body)

			fmt.Printf("%s;%v;%v;%v;%v;", s, resp.StatusCode, l, resp.ContentLength, time.Since(time_start).String())
		}
	}

}
