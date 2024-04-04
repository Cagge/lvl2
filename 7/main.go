package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	done := make(chan struct{})

	for _, channel := range channels {
		go func(ch <-chan interface{}) {

			for v := range ch {
				out <- v
			}

			done <- struct{}{}
		}(channel)
	}

	go func() {
		<-done
		close(out)
	}()

	return out
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
