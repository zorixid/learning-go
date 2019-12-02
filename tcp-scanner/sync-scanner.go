// Prints 1 -> 1024 asynchronously
// Uses a pool of goroutines, and a channel with 100 workers
// to manage the concurrent work
package main

import (
	"fmt"
	"sync"
)

// channel of type int: receives work
// pointer to a WaitGroup: tracks when work has been completed
func worker(ports chan int, wg *sync.WaitGroup) {
	// range will continuously receive from the ports channel
	// looping until the channel is closed
	for p := range ports {
		fmt.Println(p)
		wg.Done()
	}
}

func main() {
	// create a channel with `make()`. 100 workers allows the channel to be buffered
	// this means we can send an item without waiting for a receiver to read it.
	// buffered channels are useful for tracking work for multiple producers.
	// this channel will hold 100 items before the sender will block
	ports := make(chan int, 100)
	// create a synchronized counter
	var wg sync.WaitGroup
	// start the desired number of workers
	for i := 0; i < cap(ports); i++ {
		go worker(ports, &wg)
	}
	// iterate over ports sequentially
	for i := 1; i < 1024; i++ {
		// increment counter each time we create a go routine
		wg.Add(1)
		// send a port (i) on the ports channel to the worker
		ports <- i
	}
	// blocks until all work has been completed and the counter returns to zero
	wg.Wait()
	// after all the work has been complete, we'll close the channel
	close(ports)
}
