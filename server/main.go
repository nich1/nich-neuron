package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go httpListener()
	wg.Add(1)
	go timeListener()
	wg.Wait()
}
