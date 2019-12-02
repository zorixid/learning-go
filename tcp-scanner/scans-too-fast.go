// Scans ports 1 -> 65535 asynchronously. Scans too fast though.
// Network or system limitations can skew results
package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	// create a synchronized counter
	var wg sync.WaitGroup
	for i := 1; i < 65535; i++ {
		// increment counter each time we create a go routine
		wg.Add(1)
		go func(j int) {
			// decrement the counter when one unit of work has completed
			defer wg.Done()
			address := fmt.Sprintf("127.0.0.1:%d", j)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
			fmt.Printf("%d open\n", j)
		}(i)
	}
	// blocks until all work has been completed and the counter returns to zero
	wg.Wait()
}
