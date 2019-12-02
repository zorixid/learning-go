// Scans ports 1 -> 1024 using a pool of 100 workers.
// Prints the open ports to the console
package main

import (
	"fmt"
	"net"
	"sort"
)

// our worker accepts two channels as inputs
func worker(ports, results chan int) {
	// range will continuously receive input from the ports channel,
	// looping until the channel is closed
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// port is closed or filtered. Send a `0` and continue
			results <- 0
			continue
		}
		conn.Close()
		// port is open, send the port and continue
		results <- p
	}
}

func main() {
	// create a port channel with 100 workers
	ports := make(chan int, 100)
	// create a channel to communicate the results from the worker to the main thread
	results := make(chan int)
	// use a slice to store the results so we can sort them later
	var openports []int

	// start the desired number of workers
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	// send to the workers in a separate goroutine because the result-gathering loop
	// needs to start before more than 100 items of work can continue.
	// iterate over ports sequentially and send to ports channel
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// the results gathering loop receives on the results channel 1024 times
	// if the port doesn't equal 0, it's appended to the openports slice
	for i := 0; i < 1024; i++ {
		port := <-results
		fmt.Println(i)
		if port != 0 {
			openports = append(openports, port)
		}
	}

	// after all the work has been complete, we'll close the channels
	close(ports)
	close(results)

	// sort the slice of open ports
	// loop over the slice and print the open ports to the console
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
